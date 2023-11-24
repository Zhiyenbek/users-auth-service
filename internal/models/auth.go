package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserSignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Token struct {
	UserID     int64
	TokenValue string
	ExpiresAt  time.Duration
}

// Tokens - structure for holding access and refresh token
type Tokens struct {
	AccessToken  *Token
	RefreshToken *Token
}
type JwtUserClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}
