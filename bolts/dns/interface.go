package dns

import "gorm.io/gorm"

type Domain struct {
	gorm.Model
	DomainName string `gorm:"size:255;unique"`
	Provider   string
	Records    []Record `gorm:"foreignKey:DomainName;references:DomainName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Record struct {
	gorm.Model
	DomainName string `gorm:"size:255;uniqueIndex:dn_rr"`
	RR         string `gorm:"size:255;uniqueIndex:dn_rr"`
	Type       string
	Value      string
}

type Provider interface {
	AddDomain(domain *Domain) error
	HasDomain(domain *Domain) (bool, error)
	AddRecord(record *Record) error
	UpdateRecord(record *Record) error
	HasRecord(record *Record) (bool, error)
}
