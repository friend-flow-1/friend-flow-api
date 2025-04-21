package models

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	Sub string `json:"sub"` // User ID
	jwt.StandardClaims
}
