package repository

import (
	"github.com/Zhiyenbek/users-main-service/config"
	"github.com/Zhiyenbek/users-main-service/internal/models"
	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	AuthRepository
	TokenRepository
}

type AuthRepository interface {
	GetUserInfoByLogin(login string) (string, int64, error)
}

type TokenRepository interface {
	SetRTToken(token *models.Token) error
	UnsetRTToken(userID int64) error
	GetToken(userID int64) (string, error)
}

func New(db *pgxpool.Pool, cfg *config.Configs, redis *redis.Client) *Repository {
	return &Repository{
		AuthRepository:  NewAuthRepository(db, cfg.DB),
		TokenRepository: NewTokenRepository(redis),
	}
}
