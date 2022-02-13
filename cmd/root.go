package cmd

import (
	"context"
	"fmt"
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"goddns/api"
	"goddns/bolts/store"
	"goddns/config"
	"goddns/metrics"
	"os"
	"time"
)

var (
	confPath string
	debug    bool
)

var rootCmd = &cobra.Command{
	Use: "goddns",
	Run: func(c *cobra.Command, args []string) {
		svc := &svcDecorator{}
		s, err := service.New(svc, &service.Config{
			Name: "goddns",
		})
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := s.Run(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&confPath, "config", "c", "config.json", "config file")
	rootCmd.Flags().BoolVarP(&debug, "debug", "", false, "debug mode")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type svcDecorator struct{}

func (s *svcDecorator) Start(svc service.Service) error {
	conf, err := config.Load(confPath)
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}
	if debug {
		conf.API.Debug = true
		conf.API.Swagger = true
	}

	var metricsSvc metrics.Metrics
	if conf.Metrics.Enabled {
		metricsSvc, _ = metrics.NewPrometheusService(&conf.Metrics)
		go func() {
			if err := metricsSvc.Serve(); err != nil {
				log.Errorln("Metrics server error:", err)
			}
		}()
	}

	ds, err := store.NewStore(conf.ORM, metricsSvc)
	if err != nil {
		return fmt.Errorf("datastore error: %w", err)
	}

	tick, _ := time.ParseDuration(conf.DDNS.CheckInterval)
	ds.DDNSService.AutoExecute(context.TODO(), tick)

	server, err := api.NewServer(&conf.API, ds, metricsSvc)
	if err != nil {
		return fmt.Errorf("init server error: %w", err)
	}
	if err := server.Serve(); err != nil {
		return fmt.Errorf("start server error: %w", err)
	}

	return nil
}

func (s *svcDecorator) Stop(svc service.Service) error {
	return nil
}
