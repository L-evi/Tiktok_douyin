package userutil

import "github.com/golang-jwt/jwt/v4"

// JwtClaims jwt claims
type JwtClaims struct {
	ID       string `json:"jti,omitempty"`
	ExpireAt int64  `json:"exp,omitempty"`
	IssuedAt int64  `json:"iat,omitempty"`
	Issuer   string `json:"iss,omitempty"`
	Username string `json:"username"`
	UserId   int64  `json:"userId"`
	jwt.RegisteredClaims
}
