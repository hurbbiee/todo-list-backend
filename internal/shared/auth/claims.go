package auth

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID int64  `json:"userId"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}
