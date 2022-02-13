package controllers

import (
	"github.com/gin-gonic/gin"
	"goddns/bolts/store"
	"net/http"
	"strconv"
)

// @Summary Execute DDNS
// @Tags DDNS
// @Produce json
// @Param force query bool true "Force execute remotely"
// @Success 200 {object} HTTPResponse
// @Failure 400 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /ddns/execute [put]
// @Security ApiKeyAuth
func DDNSExecute(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		force, _ := strconv.ParseBool(c.Query("force"))
		if err := store.DDNSService.Execute(force); err != nil {
			c.JSON(http.StatusInternalServerError, HTTPResponse{
				Code: http.StatusInternalServerError, Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, HTTPResponse{Code: http.StatusOK})
	}
}
