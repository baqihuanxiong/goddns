package controllers

import (
	"github.com/gin-gonic/gin"
	"goddns/bolts/crypto"
	"goddns/bolts/store"
	"goddns/bolts/user"
	"net/http"
)

type authPayload struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type authResponse struct {
	JWT string `json:"jwt"`
}

// @Summary Request an access token
// @Tags Authenticate
// @Accept json
// @Produce json
// @Param payload body authPayload true "Authentication payload"
// @Success 200 {object} authResponse
// @Failure 400 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /login [post]
func Authenticate(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// support internal user authentication
		authenticateInternal(c, store)
	}
}

func authenticateInternal(c *gin.Context, store *store.Store) {
	var payload authPayload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, HTTPResponse{
			Code: http.StatusBadRequest, Error: "invalid request body: " + err.Error()})
		return
	}

	ok, err := store.UserService.MatchUserWithPlainPassword(&user.User{
		Username: payload.Username,
		Password: payload.Password,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, HTTPResponse{
			Code: http.StatusForbidden, Error: "invalid user"})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, HTTPResponse{
			Code: http.StatusForbidden, Error: "invalid password"})
		return
	}

	// create jwt token for authenticated user
	j := crypto.NewJWT()
	token, err := j.CreateToken(crypto.CustomClaims{
		Username: payload.Username,
	})
	c.JSON(http.StatusOK, authResponse{JWT: token})
}
