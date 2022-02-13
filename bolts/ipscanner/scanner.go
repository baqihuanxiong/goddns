package ipscanner

import (
	"errors"
	"goddns/config"
	"gorm.io/gorm"
)

var (
	ErrNoAvailableProvider = errors.New("no scanner available")
)

type Service struct {
	db *gorm.DB

	providers []Provider
	cur       int

	// current dynamic ip
	current string
}

func NewIPScannerService(db *gorm.DB, autoMigrate bool, providers []config.IPScanner) (*Service, error) {
	svc := &Service{
		db:        db,
		providers: make([]Provider, 0),
	}
	if autoMigrate {
		if err := db.AutoMigrate(&Status{}); err != nil {
			return svc, err
		}
	}
	if lastStatus, err := svc.LastStatus(); err == nil {
		svc.current = lastStatus.Current
	}
	return svc, svc.initProviders(providers)
}

func (s *Service) initProviders(providers []config.IPScanner) error {
	for _, provider := range providers {
		if provider.Provider == "customScanner" {
			p, err := NewCustomScanner(provider.CustomScanner)
			if err != nil {
				return err
			}
			s.providers = append(s.providers, p)
		} else if provider.Provider == "jsonip" {
			p, err := NewJSONIPScanner(provider.JSONIPScanner)
			if err != nil {
				return err
			}
			s.providers = append(s.providers, p)
		}
	}
	return nil
}

func (s *Service) LastStatus() (Status, error) {
	var status Status
	tx := s.db.Last(&status)
	return status, tx.Error
}

// Current dynamic ip
func (s *Service) Current() string {
	return s.current
}

// Lookup return client's public IP address, it will check whether
// current provider is alive, if not it will try next provider
// in ipScanners list, if no provider is alive after two
// round of check, it will return ErrProviderUnhealthy error
func (s *Service) Lookup() (string, error) {
	for i := 0; i < len(s.providers)*2; i++ {
		provider := s.providers[s.cur]
		ip, err := provider.Lookup()
		if err != nil {
			s.cur = (s.cur + 1) % len(s.providers)
			continue
		}
		return ip, s.saveStatus(ip)
	}
	return "", ErrNoAvailableProvider
}

func (s *Service) saveStatus(ip string) error {
	if s.current == ip {
		return nil
	}
	s.current = ip
	return s.db.Create(&Status{Current: ip}).Error
}
