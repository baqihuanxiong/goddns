package cmd

import (
	"context"
	"fmt"
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
	logLevel int32
)

var rootCmd = &cobra.Command{
	Use: "goddns",
	Run: func(c *cobra.Command, args []string) {
		conf, err := config.Load(confPath)
		if err != nil {
			log.Fatalln("config error: %w", err)
		}
		if debug {
			conf.API.Debug = true
			conf.API.Swagger = true
		}
		log.Infoln("loaded config:", conf.String())

		var metricsSvc metrics.Metrics
		if conf.Metrics.Enabled {
			metricsSvc, err = metrics.NewPrometheusService(&conf.Metrics)
			if err != nil {
				log.Fatalln("initialize metrics server error:", err)
			}
			log.Infoln("metrics server initialized")
			go func() {
				if err := metricsSvc.Serve(); err != nil {
					log.Fatalln("metrics server error:", err)
				}
			}()
		}

		ds, err := store.NewStore(conf.ORM, metricsSvc)
		if err != nil {
			log.Fatalln("datastore error: %w", err)
		}
		log.Infoln("datastore established")

		tick, _ := time.ParseDuration(conf.DDNS.CheckInterval)
		ds.DDNSService.AutoExecute(context.TODO(), tick)

		server, err := api.NewServer(&conf.API, ds, metricsSvc)
		if err != nil {
			log.Fatalln("init server error: %w", err)
		}
		if err := server.Serve(); err != nil {
			log.Fatalln("start server error: %w", err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&confPath, "config", "c", "config.json", "config file")
	rootCmd.Flags().BoolVarP(&debug, "debug", "", false, "debug mode")
	rootCmd.Flags().Int32VarP(&logLevel, "log-level", "l", 3, "log level, 1-5")
	log.SetLevel(log.Level(logLevel))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
