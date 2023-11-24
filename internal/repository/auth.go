package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zhiyenbek/users-main-service/config"
	"github.com/Zhiyenbek/users-main-service/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type authRepository struct {
	db  *pgxpool.Pool
	cfg *config.DBConf
}

func NewAuthRepository(db *pgxpool.Pool, cfg *config.DBConf) AuthRepository {
	return &authRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *authRepository) GetUserInfoByLogin(login string) (string, int64, error) {
	var password string
	var ID int64
	timeout := r.cfg.TimeOut
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	query := "SELECT password, user_id FROM auth WHERE login = $1"
	if err := r.db.QueryRow(ctx, query, login).Scan(&password, &ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", -1, fmt.Errorf("%w: error occurred while getting password_hash from db: auth does not exist", models.ErrWrongPassword)
		}
		return "", -1, fmt.Errorf("%w: error occurred while getting password_hash from db: %v", models.ErrInternalServer, err)
	}
	return password, ID, nil

}
