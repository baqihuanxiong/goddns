package ddns

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"goddns/bolts/dns"
	"goddns/bolts/ipscanner"
	"goddns/metrics"
	"time"
)

type Service struct {
	dnsService     *dns.Service
	scannerService *ipscanner.Service
	metrics        metrics.Metrics
}

func NewDDNSService(dnsService *dns.Service, scannerService *ipscanner.Service,
	metricsService metrics.Metrics) (*Service, error) {
	svc := &Service{
		dnsService:     dnsService,
		scannerService: scannerService,
		metrics:        metricsService,
	}
	return svc, nil
}

func (s *Service) ExecuteOne(domainName, rr, dnsType, value string) error {
	domain, err := s.dnsService.GetDomainByName(domainName)
	if err != nil {
		return err
	}
	provider, ok := s.dnsService.GetProvider(domain.Provider)
	if !ok {
		return dns.ErrNoProvider
	}
	if err := provider.UpdateRecord(&dns.Record{
		DomainName: domainName,
		RR:         rr,
		Type:       dnsType,
		Value:      value,
	}); err != nil {
		return err
	}

	if s.metrics != nil {
		s.metrics.IncreaseDDNSCounter()
	}
	return nil
}

func (s *Service) lookupIP() (string, bool, error) {
	currentIP := s.scannerService.Current()
	ip, err := s.scannerService.Lookup()
	if err != nil {
		return "", false, err
	}
	if ip == currentIP {
		return ip, false, nil
	}
	return ip, true, nil
}

// Execute DDNS operation in batch, if forceUpdateRemote is true, remote dns will be executed
// regardless whether dynamic IP has changed.
func (s *Service) Execute(forceUpdateRemote bool) error {
	ip, isNew, err := s.lookupIP()
	if err != nil {
		return err
	}
	log.Debugln("lookup ip:", ip)
	if !forceUpdateRemote && !isNew {
		return nil
	}

	// execute dns only on records with no value set
	log.Infoln("DDNS service detected new ip, updating remote...")
	records, err := s.dnsService.ListRecords()
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Value != "" {
			continue
		}
		if err := s.ExecuteOne(record.DomainName, record.RR, record.Type, ip); err != nil {
			log.Errorln(fmt.Sprintf(
				"DDNS service update record %s error: %v",
				record.RR+"."+record.DomainName, err))
			continue
		}
		log.Warningln(fmt.Sprintf(
			"DDNS service updated remote record %s value to %s",
			record.RR+"."+record.DomainName, ip))
	}
	return nil
}

func (s *Service) AutoExecute(ctx context.Context, tick time.Duration) {
	if err := s.Execute(true); err != nil {
		log.Errorln("DDNS service execute error:", err)
	}

	go func() {
		ticker := time.NewTicker(tick)
	loop:
		for {
			select {
			case <-ticker.C:
				if err := s.Execute(false); err != nil {
					log.Errorln("DDNS service execute error:", err)
				}
			case <-ctx.Done():
				break loop
			}
		}
	}()
}
