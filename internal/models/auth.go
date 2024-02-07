package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type UserSignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserData struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RecruiterSignUpRequest struct {
	UserData
	CompanyName     string    `json:"company_name"`
	CompanyPublicID uuid.UUID `json:"company_public_id"`
}

type CandidateSignUpRequest struct {
	UserData
	Resume          string   `json:"resume"`
	CurrentPosition string   `json:"current_position"`
	Bio             string   `json:"bio"`
	Skills          []string `json:"skills"`
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
