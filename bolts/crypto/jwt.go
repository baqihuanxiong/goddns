package crypto

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var (
	ErrMalformedToken    = errors.New("malformed token")
	ErrTokenExpired      = errors.New("expired token")
	ErrNotValidatedToken = errors.New("token not validated")
	ErrUnknown           = errors.New("unknown token")
	ErrInvalidHeader     = errors.New("invalid token header")
)

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte("7a6a3a"),
	}
}

type CustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	if len(tokenString) <= 7 || tokenString[:7] != "Bearer " {
		return nil, ErrInvalidHeader
	}
	token, err := jwt.ParseWithClaims(tokenString[7:], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrMalformedToken
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrNotValidatedToken
			}
		}
		return nil, ErrUnknown
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrUnknown
}
