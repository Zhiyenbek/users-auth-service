package service

import (
	"fmt"
	"time"

	"github.com/Zhiyenbek/users-auth-service/config"
	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/Zhiyenbek/users-auth-service/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	cfg           *config.Configs
	logger        *zap.SugaredLogger
	authRepo      repository.AuthRepository
	tokenRepo     repository.TokenRepository
	recruiterRepo repository.RecruiterRepository
	candidateRepo repository.CandidateRepository
}

func NewAuthService(repo *repository.Repository, cfg *config.Configs, logger *zap.SugaredLogger) AuthService {
	return &authService{
		authRepo:      repo.AuthRepository,
		tokenRepo:     repo.TokenRepository,
		recruiterRepo: repo.RecruiterRepository,
		candidateRepo: repo.CandidateRepository,
		cfg:           cfg,
		logger:        logger,
	}
}

func (s *authService) CreateCandidate(req *models.CandidateSignUpRequest) error {
	var err error
	req.Password, err = hashAndSalt([]byte(req.Password))
	if err != nil {
		s.logger.Error("could not hash password")
		return err
	}
	err = s.candidateRepo.CreateCandidate(req)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) CreateRecruiter(req *models.RecruiterSignUpRequest) error {
	var err error
	req.Password, err = hashAndSalt([]byte(req.Password))
	if err != nil {
		s.logger.Error("could not hash password")
		return err
	}
	err = s.recruiterRepo.CreateRecruiter(req)
	if err != nil {
		return err
	}
	return nil
}
func (s *authService) Login(creds *models.UserSignInRequest) (*models.Tokens, error) {
	pass, userID, err := s.authRepo.GetUserInfoByLogin(creds.Login)
	if err != nil {
		return nil, err
	}
	if !checkPasswordHash(creds.Password, pass) {
		s.logger.Error("failed to login. Password didn't match")
		return nil, models.ErrWrongPassword
	}

	return s.generateTokens(userID)
}

// hashAndSalt - hashes the password with salt. Function takes password as []byte and returns the hash as string and error.
func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash - checks if the password matches the hash. Function takes password and has as string, and returns true if they matched and false otherwise.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateAccessToken - function for creating new access token for user
func createAccessToken(userID int64, tokenTTL time.Duration, tokenSecret string) (*models.Token, error) {
	var err error
	//Creating Access Token
	iat := time.Now().Unix()
	exp := time.Now().Add(tokenTTL)
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["iat"] = iat
	atClaims["exp"] = exp.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	tokenString, err := at.SignedString([]byte(tokenSecret))
	if err != nil {
		return nil, err
	}
	token := &models.Token{
		TokenValue: tokenString,
		UserID:     userID,
		ExpiresAt:  time.Until(exp),
	}
	return token, nil
}

// CreateRefreshToken - function for creating new refresh token for user
func createRefreshToken(userID int64, tokenTTL time.Duration, tokenSecret string) (*models.Token, error) {
	var err error
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	iat := time.Now().Unix()
	exp := time.Now().Add(tokenTTL)
	rtClaims["authorized"] = true
	rtClaims["user_id"] = userID
	rtClaims["iat"] = iat
	rtClaims["exp"] = exp.Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	tokenString, err := at.SignedString([]byte(tokenSecret))
	if err != nil {
		return nil, err
	}
	token := &models.Token{
		TokenValue: tokenString,
		UserID:     userID,
		ExpiresAt:  time.Until(exp),
	}
	return token, nil
}

// ParseToken - method that responsible for parsing jwt token. It checks if jwt token is valid, retrieves claims and returns user public id. In case of error returns error
func (s *authService) parseToken(tokenString string, tokenSecret string) (*models.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		s.logger.Error(err)
		return nil, fmt.Errorf("could not parse token: %v %w", err, models.ErrInvalidToken)
	}
	if claims, ok := token.Claims.(*models.JwtUserClaims); ok && token.Valid {
		token := &models.Token{
			UserID:     claims.UserID,
			TokenValue: tokenString,
		}
		return token, nil
	}
	return nil, fmt.Errorf("could not parse token: %w", models.ErrInvalidToken)
}

func (s *authService) RefreshToken(tokenString string) (*models.Tokens, error) {
	token, err := s.parseToken(tokenString, s.cfg.Token.Refresh.TokenSecret)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	redisTokenString, err := s.tokenRepo.GetToken(token.UserID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	if tokenString != redisTokenString {
		s.logger.Errorf("token is unmatched. Wanted %s. Got: %s", tokenString, redisTokenString)
		return nil, models.ErrTokenExpired
	}
	err = s.tokenRepo.UnsetRTToken(token.UserID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	tokens, err := s.generateTokens(token.UserID)

	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return tokens, nil
}

// GenerateTokens - method that responsible for generating tokens. It generates jwt access token and refresh token and returns them as models.Tokenss. In case of error returns error
func (s *authService) generateTokens(userID int64) (*models.Tokens, error) {
	accessToken, err := createAccessToken(userID, s.cfg.Token.Access.ExpiresAt, s.cfg.Token.Access.TokenSecret)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	refreshToken, err := createRefreshToken(userID, s.cfg.Token.Refresh.ExpiresAt, s.cfg.Token.Refresh.TokenSecret)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	err = s.tokenRepo.SetRTToken(refreshToken)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	tokens := &models.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}
	return tokens, nil
}
