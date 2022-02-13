package controllers

import (
	"github.com/gin-gonic/gin"
	"goddns/bolts/dns"
	"goddns/bolts/store"
	"net/http"
)

type Domain struct {
	DomainName string `form:"dn" json:"dn" binding:"required"`
	Provider   string `form:"provider" json:"provider" binding:"required"`
}

type Record struct {
	DomainName string `form:"dn" json:"dn" binding:"required"`
	RR         string `form:"rr" json:"rr" binding:"required"`
	Type       string `form:"type" json:"type" binding:"required,oneof=A AAAA NS CNAME"`
	Value      string `form:"value" json:"value"`
}

func transDomain(domain *dns.Domain) *Domain {
	return &Domain{
		DomainName: domain.DomainName,
		Provider:   domain.Provider,
	}
}

func transRecord(record *dns.Record) *Record {
	return &Record{
		DomainName: record.DomainName,
		RR:         record.RR,
		Type:       record.Type,
		Value:      record.Value,
	}
}

// @Summary Add a domain
// @Tags DNS
// @Accept json
// @Produce json
// @Param payload body Domain true "Domain"
// @Success 200 {object} Domain
// @Failure 400 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /domain [post]
// @Security ApiKeyAuth
func AddDomain(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload Domain
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, HTTPResponse{
				Code: http.StatusBadRequest, Error: err.Error()})
			return
		}
		if err := store.DNSService.AddDomain(&dns.Domain{
			DomainName: payload.DomainName,
			Provider:   payload.Provider,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, HTTPResponse{
				Code: http.StatusInternalServerError, Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}

// @Summary Delete a domain
// @Tags DNS
// @Produce json
// @Param dn query string true "Domain name"
// @Success 200 {object} HTTPResponse
// @Failure 400 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /domain [delete]
// @Security ApiKeyAuth
func DeleteDomain(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		dn := c.Query("dn")
		if err := store.DNSService.DeleteDomain(&dns.Domain{DomainName: dn}); err != nil {
			c.JSON(http.StatusInternalServerError, HTTPResponse{
				Code: http.StatusInternalServerError, Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, HTTPResponse{Code: http.StatusOK})
	}
}

// @Summary List domain providers
// @Tags DNS
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} HTTPResponse
// @Router /domain/provider/list [get]
// @Security ApiKeyAuth
func ListDomainProviders(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		providers := store.DNSService.ListProviders()
		c.JSON(http.StatusOK, providers)
	}
}

// @Summary List domains
// @Tags DNS
// @Produce json
// @Success 200 {array} Domain
// @Failure 500 {object} HTTPResponse
// @Router /domain/list [get]
// @Security ApiKeyAuth
func ListDomains(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		domains, err := store.DNSService.ListDomains()
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPResponse{
				Code: http.StatusInternalServerError, Error: err.Error()})
			return
		}
		transDomains := make([]*Domain, 0, len(domains))
		for _, domain := range domains {
			transDomains = append(transDomains, transDomain(domain))
		}
		c.JSON(http.StatusOK, transDomains)
	}
}

// @Summary Add a DNS record
// @Tags DNS
// @Accept json
// @Produce json
// @Param payload body Record true "DNS record"
// @Success 200 {object} Record
// @Failure 400 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /record [post]
// @Security ApiKeyAuth
func AddRecord(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload Record
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, HTTPResponse{
				Code: http.StatusBadRequest, Error: err.Error()})
			return
		}
		if err := store.DNSService.AddRecord(&dns.Record{
			DomainName: payload.DomainName,
			RR:         payload.RR,
			Type:       payload.Type,
			Value:      payload.Value,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, HTTPResponse{
				Code: http.StatusInternalServerError, Error: err.Error()})
			return
		}
		if payload.Value == "" {
			ip := store.IPScanService.Current()
			if err := store.DDNSService.ExecuteOne(payload.DomainName, payload.RR, payload.Type, ip); err != nil {
				c.JSON(http.StatusInternalServerError, HTTPResponse{
					Code: http.StatusInternalServerError, Error: err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, payload)
	}
}

// @Summary Delete a DNS record
// @Tags DNS
// @Produce json
// @Param dn query string true "Domain of record"
// @Param rr query string true "Record resource"
// @Success 200 {object} HTTPResponse
// @Failure 400 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /record [delete]
// @Security ApiKeyAuth
func DeleteRecord(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		dn := c.Query("dn")
		rr := c.Query("rr")
		if err := store.DNSService.DeleteRecord(&dns.Record{DomainName: dn, RR: rr}); err != nil {
			c.JSON(http.StatusInternalServerError, HTTPResponse{
				Code: http.StatusInternalServerError, Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, HTTPResponse{Code: http.StatusOK})
	}
}

// @Summary List DNS records
// @Tags DNS
// @Produce json
// @Param dn query string true "Domain name"
// @Success 200 {array} Record
// @Failure 500 {object} HTTPResponse
// @Router /record/list [get]
// @Security ApiKeyAuth
func ListRecords(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		dn := c.Query("dn")
		var records []*dns.Record
		var err error
		if dn == "" {
			records, err = store.DNSService.ListRecords()
		} else {
			records, err = store.DNSService.ListRecordsByDomain(dn)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPResponse{
				Code: http.StatusInternalServerError, Error: err.Error()})
			return
		}
		transRecords := make([]*Record, 0, len(records))
		for _, record := range records {
			transRecords = append(transRecords, transRecord(record))
		}
		c.JSON(http.StatusOK, transRecords)
	}
}
