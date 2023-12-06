package domain

import (
	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
)

type JwtCustomClaims struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
	jwt.StandardClaims
}

type JwtCustomRefreshClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}
