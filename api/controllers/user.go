package controllers

import (
	"github.com/gin-gonic/gin"
	"goddns/bolts/crypto"
	"goddns/bolts/store"
	"net/http"
)

type User struct {
	Username string `json:"username"`
}

// @Summary Current authenticated user
// @Tags User
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /user/current [get]
// @Security ApiKeyAuth
func ParseUser(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		iClamis, ok := c.Get("claims")
		if !ok {
			c.JSON(http.StatusBadRequest, HTTPResponse{
				Code: http.StatusBadRequest, Message: "authentication not found"})
			return
		}
		claims, ok := iClamis.(*crypto.CustomClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, HTTPResponse{
				Code: http.StatusInternalServerError, Message: "invalid authentication"})
			return
		}
		c.JSON(http.StatusOK, &User{Username: claims.Username})
	}
}
