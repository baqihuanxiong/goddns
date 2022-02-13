package dns

import (
	"errors"
	"goddns/config"
	"gorm.io/gorm"
)

var (
	ErrNoProvider = errors.New("no dns provider")
)

type Service struct {
	db        *gorm.DB
	providers map[string]Provider
}

func NewDNSService(db *gorm.DB, autoMigrate bool, providers []config.DNS) (*Service, error) {
	svc := &Service{
		db:        db,
		providers: make(map[string]Provider),
	}
	if autoMigrate {
		if err := db.AutoMigrate(&Domain{}, &Record{}); err != nil {
			return nil, err
		}
	}
	return svc, svc.initProviders(providers)
}

func (s *Service) initProviders(providers []config.DNS) error {
	for _, provider := range providers {
		var err error
		if provider.Provider == "alidns" {
			if s.providers[provider.Provider], err = NewAliDNSService(provider.AliDNS); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) ListProviders() []string {
	providers := make([]string, 0)
	for name, _ := range s.providers {
		providers = append(providers, name)
	}
	return providers
}

func (s *Service) GetProvider(name string) (Provider, bool) {
	provider, ok := s.providers[name]
	return provider, ok
}

func (s *Service) AddDomain(domain *Domain) error {
	provider, ok := s.providers[domain.Provider]
	if !ok {
		return ErrNoProvider
	}
	hasDomain, err := provider.HasDomain(domain)
	if err != nil {
		return err
	}
	if !hasDomain {
		if err := provider.AddDomain(domain); err != nil {
			return err
		}
	}

	return s.db.Create(domain).Error
}

func (s *Service) DeleteDomain(domain *Domain) error {
	return s.db.Unscoped().Delete(&Domain{}, "domain_name = ?", domain.DomainName).Error
}

func (s *Service) GetDomainByName(domainName string) (*Domain, error) {
	var domain Domain
	tx := s.db.First(&domain, "domain_name = ?", domainName)
	return &domain, tx.Error
}

func (s *Service) ListDomains() ([]*Domain, error) {
	domains := make([]*Domain, 0)
	tx := s.db.Find(&domains)
	return domains, tx.Error
}

func (s *Service) AddRecord(record *Record) error {
	domain, err := s.GetDomainByName(record.DomainName)
	if err != nil {
		return err
	}
	provider, ok := s.providers[domain.Provider]
	if !ok {
		return ErrNoProvider
	}
	hasRecord, err := provider.HasRecord(record)
	if err != nil {
		return err
	}
	if !hasRecord {
		if err := provider.AddRecord(record); err != nil {
			return err
		}
	}

	return s.db.Create(record).Error
}

func (s *Service) DeleteRecord(record *Record) error {
	return s.db.Unscoped().Delete(&Record{}, "domain_name = ? and rr = ?", record.DomainName, record.RR).Error
}

func (s *Service) ListRecords() ([]*Record, error) {
	records := make([]*Record, 0)
	tx := s.db.Find(&records)
	return records, tx.Error
}

func (s *Service) ListRecordsByDomain(domainName string) ([]*Record, error) {
	records := make([]*Record, 0)
	tx := s.db.Find(&records, "domain_name = ?", domainName)
	return records, tx.Error
}
