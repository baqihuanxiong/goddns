package ipscanner

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"goddns/config"
	"io/ioutil"
	"net/http"
	"time"
)

type CustomScanner struct {
	conf    *config.CustomScanner
	httpCli *http.Client
}

type IPResponse struct {
	Status  int    `json:"status"`
	IP      string `json:"ip"`
	Message string `json:"message,omitempty"`
}

func NewCustomScanner(config *config.CustomScanner) (*CustomScanner, error) {
	httpCli := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	return &CustomScanner{conf: config, httpCli: httpCli}, nil
}

func (s *CustomScanner) Lookup() (string, error) {
	url := s.conf.URL + "?token=" + s.conf.Token
	resp, err := s.httpCli.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	buf, _ := ioutil.ReadAll(resp.Body)
	var ip IPResponse
	if err := json.Unmarshal(buf, &ip); err != nil {
		return "", err
	}
	if ip.Status != 0 {
		return "", errors.New(ip.Message)
	}
	return ip.IP, nil
}
