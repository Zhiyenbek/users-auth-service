package service

import (
	"github.com/Zhiyenbek/users-main-service/config"
	"github.com/Zhiyenbek/users-main-service/internal/models"
	"github.com/Zhiyenbek/users-main-service/internal/repository"
	"go.uber.org/zap"
)

type AuthService interface {
	Login(creds *models.UserSignInRequest) (*models.Tokens, error)
	RefreshToken(tokenString string) (*models.Tokens, error)
}

type Service struct {
	AuthService
}

func New(repos *repository.Repository, log *zap.SugaredLogger, cfg *config.Configs) *Service {
	return &Service{
		AuthService: NewAuthService(repos, cfg, log),
	}
}
