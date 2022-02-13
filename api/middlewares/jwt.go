package middlewares

import (
	"github.com/gin-gonic/gin"
	"goddns/bolts/crypto"
	"net/http"
)

// JWTWsAuthTrans is inserted in front of JWTAuth to protect websocket connect as well
func JWTWsAuthTrans() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Sec-WebSocket-Protocol")
		c.Request.Header.Set("Authorization", "Bearer "+token)
	}
}

// JWTAuth provide the most basic authentication to access API, reference: https://jwt.io/introduction
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "require authentication token",
			})
			c.Abort()
			return
		}

		j := crypto.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		// write validated claims back to gin.Context
		c.Set("claims", claims)
	}
}
