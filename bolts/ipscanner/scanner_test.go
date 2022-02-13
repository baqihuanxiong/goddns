package ipscanner

import (
	"goddns/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func newSvc() (*Service, error) {
	db, err := gorm.Open(sqlite.Open("goddns.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return NewIPScannerService(db, true, []config.IPScanner{
		{
			Provider: "jsonip",
			JSONIPScanner: &config.JSONIPScanner{
				URL: "https://jsonip.com",
			},
		},
	})
}

func TestService_Lookup(t *testing.T) {
	svc, err := newSvc()
	if err != nil {
		t.Fatal(err)
	}

	ip, err := svc.Lookup()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip)
}
