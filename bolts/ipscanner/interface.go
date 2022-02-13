package ipscanner

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	Current string
}

type Provider interface {
	// Lookup current ip address by querying third-party api or reading from local system
	Lookup() (string, error)
}
