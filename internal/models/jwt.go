package models

import "github.com/golang-jwt/jwt/v4"

type JWTClaims struct {
	Sub string `json:"sub"` // User ID
	jwt.StandardClaims
}
