package store

import (
	"errors"
	"goddns/bolts/ddns"
	"goddns/bolts/dns"
	"goddns/bolts/ipscanner"
	"goddns/bolts/user"
	"goddns/config"
	"goddns/metrics"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrUnsupportedDB = errors.New("unsupported db driver")
)

type Store struct {
	conf    config.ORM
	db      *gorm.DB
	metrics metrics.Metrics

	UserService   *user.Service
	DNSService    *dns.Service
	IPScanService *ipscanner.Service
	DDNSService   *ddns.Service
}

// Store provider internal db storage and other services,
// it's supposed to be injected into API server service
func NewStore(config config.ORM, metricsService metrics.Metrics) (*Store, error) {
	s := &Store{conf: config, metrics: metricsService}
	if err := s.open(); err != nil {
		return nil, err
	}
	return s, s.initServices()
}

// connect to db, db file will be created if not exists when using sqlite
func (s *Store) open() error {
	var err error
	if s.conf.Driver == "sqlite" {
		s.db, err = gorm.Open(sqlite.Open(s.conf.DSN), &gorm.Config{})
		if err != nil {
			return err
		}
		// SQLite defaults to PRAGMA foreign_keys = OFF every time you open the database.
		// It's not a property of a table or of the schema.
		s.db.Exec("PRAGMA foreign_keys = ON")
	} else if s.conf.Driver == "mysql" {
		s.db, err = gorm.Open(mysql.Open(s.conf.DSN), &gorm.Config{})
		if err != nil {
			return err
		}
	} else {
		return ErrUnsupportedDB
	}
	return nil
}

func (s *Store) initServices() error {
	var err error
	if s.UserService, err = user.NewUserService(s.db, s.conf.AutoMigrate); err != nil {
		return err
	}
	if err := s.createDefaultUser(); err != nil {
		return err
	}
	if s.DNSService, err = dns.NewDNSService(s.db, s.conf.AutoMigrate, s.conf.DNSProviders); err != nil {
		return err
	}
	if s.IPScanService, err = ipscanner.NewIPScannerService(s.db, s.conf.AutoMigrate, s.conf.IPScanners); err != nil {
		return err
	}
	if s.DDNSService, err = ddns.NewDDNSService(s.DNSService, s.IPScanService, s.metrics); err != nil {
		return err
	}
	return nil
}

func (s *Store) createDefaultUser() error {
	if _, err := s.UserService.GetUserByUsername("admin"); err != nil {
		return s.UserService.CreateUser(&user.User{
			Username: "admin",
			Password: s.conf.AdminPassword,
		})
	}
	return nil
}
