package ipscanner

import (
	"crypto/tls"
	"encoding/json"
	"goddns/config"
	"io/ioutil"
	"net/http"
	"time"
)

type JSONIPScanner struct {
	conf    *config.JSONIPScanner
	httpCli *http.Client
}

type JSONIPResponse struct {
	IP string `json:"ip"`
}

func NewJSONIPScanner(config *config.JSONIPScanner) (*JSONIPScanner, error) {
	httpCli := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	return &JSONIPScanner{conf: config, httpCli: httpCli}, nil
}

func (s *JSONIPScanner) Lookup() (string, error) {
	url := s.conf.URL
	resp, err := s.httpCli.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	buf, _ := ioutil.ReadAll(resp.Body)
	var ip JSONIPResponse
	if err := json.Unmarshal(buf, &ip); err != nil {
		return "", err
	}
	return ip.IP, nil
}
