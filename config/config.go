package config

import (
	"encoding/json"
	"errors"
	"goddns/bolts/crypto"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	DDNS    DDNS    `json:"ddns"`
	API     API     `json:"api"`
	ORM     ORM     `json:"orm"`
	Metrics Metrics `json:"metrics"`
}

type DDNS struct {
	CheckInterval string `json:"checkInterval"`
}

type IPScanner struct {
	Provider      string         `json:"provider"`
	CustomScanner *CustomScanner `json:"customScanner,omitempty"`
	JSONIPScanner *JSONIPScanner `json:"jsonipScanner,omitempty"`
}

type CustomScanner struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type JSONIPScanner struct {
	URL string `json:"url"`
}

type DNS struct {
	Provider string  `json:"provider"`
	AliDNS   *AliDNS `json:"alidns,omitempty"`
}

type AliDNS struct {
	Region          string `json:"region"`
	AccessKey       string `json:"accessKey"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type API struct {
	Debug             bool   `json:"debug"`
	Listen            string `json:"listen"`
	Swagger           bool   `json:"swagger"`
	AllowCors         bool   `json:"allowCors"`
	AllowInsecureHTTP bool   `json:"allowInsecureHttp"`
	TLSKey            string `json:"tlsKey"`
	TLSCert           string `json:"tlsCert"`
	UIResource        string `json:"uiResource"`
	UIPath            string `json:"uiPath"`
}

type ORM struct {
	Driver        string      `json:"driver"`
	DSN           string      `json:"dsn"`
	AutoMigrate   bool        `json:"autoMigrate"`
	AdminPassword string      `json:"adminPassword"`
	IPScanners    []IPScanner `json:"ipScanners"`
	DNSProviders  []DNS       `json:"dnsProviders"`
}

type Metrics struct {
	Enabled         bool   `json:"enabled"`
	Listen          string `json:"listen"`
	MetricsPath     string `json:"metricsPath"`
	HealthCheckPath string `json:"healthCheckPath"`
}

func Load(path string) (*Config, error) {
	conf := &Config{}
	conf.LoadDefault()
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if _, err := conf.Verify(); err != nil {
			return nil, err
		}
		return conf.WriteToFile(path)
	}
	if _, err := conf.LoadFromFile(path); err != nil {
		return nil, err
	}
	return conf.Verify()
}

func (c *Config) LoadDefault() *Config {
	c.DDNS = DDNS{
		CheckInterval: "10s",
	}
	c.API = API{
		Listen:            ":8001",
		AllowCors:         true,
		AllowInsecureHTTP: false,
		UIResource:        "dist",
		UIPath:            "/",
	}
	c.ORM = ORM{
		Driver:        "sqlite",
		DSN:           "goddns.db",
		AutoMigrate:   true,
		AdminPassword: "admin",
		IPScanners: []IPScanner{
			{
				Provider: "jsonip",
				JSONIPScanner: &JSONIPScanner{
					URL: "https://jsonip.com",
				},
			},
		},
		DNSProviders: []DNS{
			{
				Provider: "alidns",
				AliDNS:   &AliDNS{},
			},
		},
	}
	c.Metrics = Metrics{
		Enabled:         true,
		Listen:          ":9002",
		MetricsPath:     "/metrics",
		HealthCheckPath: "/ping",
	}
	return c
}

func (c *Config) Verify() (*Config, error) {
	checkInterval, err := time.ParseDuration(c.DDNS.CheckInterval)
	if err != nil {
		return c, errors.New("invalid ddns.checkInterval")
	}
	if checkInterval < time.Second {
		return c, errors.New("ddns.checkInterval must be longer than 1s")
	}
	if len(c.ORM.IPScanners) == 0 {
		return c, errors.New("empty ipScanners")
	}
	if len(c.ORM.DNSProviders) == 0 {
		return c, errors.New("empty dnsProviders")
	}
	if !c.API.AllowInsecureHTTP && (c.API.TLSCert == "" || c.API.TLSKey == "") {
		_ = os.Mkdir("certs", 0644)
		if err := crypto.GenerateCertAndKey("certs/server.crt", "certs/server.key"); err != nil {
			return nil, err
		}
		c.API.TLSCert = "certs/server.crt"
		c.API.TLSKey = "certs/server.key"
	}
	return c, nil
}

func (c *Config) String() string {
	b, err := json.MarshalIndent(c, "", "  ")
	if err == nil {
		return string(b)
	}
	return ""
}

func (c *Config) LoadFromFile(path string) (*Config, error) {
	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(cfg, c)
	return c, err
}

func (c *Config) WriteToFile(path string) (*Config, error) {
	err := ioutil.WriteFile(path, []byte(c.String()), 0644)
	return c, err
}
