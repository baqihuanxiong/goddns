package controllers

import (
	"github.com/gin-gonic/gin"
	"goddns/bolts/store"
	"net/http"
	"time"
)

type IPStatus struct {
	IP        string `json:"ip"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// @Summary Lookup dynamic IP
// @Tags IP
// @Produce json
// @Success 200 {array} IPStatus
// @Failure 500 {object} HTTPResponse
// @Router /ip/lookup [get]
// @Security ApiKeyAuth
func LookupCurrentIP(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, err := store.IPScanService.Lookup()
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPResponse{
				Code: http.StatusInternalServerError, Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, &IPStatus{IP: ip})
	}
}

// @Summary Last dynamic IP
// @Tags IP
// @Produce json
// @Success 200 {array} IPStatus
// @Failure 500 {object} HTTPResponse
// @Router /ip/last [get]
// @Security ApiKeyAuth
func GetLastIP(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		status, err := store.IPScanService.LastStatus()
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPResponse{
				Code: http.StatusInternalServerError, Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, &IPStatus{
			IP:        status.Current,
			UpdatedAt: status.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}
}
