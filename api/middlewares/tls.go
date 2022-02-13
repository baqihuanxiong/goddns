package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func TlsHandler(opts secure.Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(opts)
		if err := secureMiddleware.Process(c.Writer, c.Request); err != nil {
			return
		}

		c.Next()
	}
}
