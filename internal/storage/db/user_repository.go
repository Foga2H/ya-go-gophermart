package db

import (
	"context"
	"database/sql"

	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/model"
	"github.com/google/uuid"
)

type UserRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewUserRepository(db *sql.DB, logger *logger.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
}

func (s *UserRepository) Create(ctx context.Context, login string, hashedPassword string) (model.User, error) {
	var user model.User
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO users (id, login, password) VALUES ($1, $2, $3) RETURNING id, login, password, created_at`,
		uuid.NewString(),
		login,
		hashedPassword).Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)

	if err != nil {
		s.logger.Errorf("failed to insert user: %v", err)
		return model.User{}, err
	}

	return user, nil
}

func (s *UserRepository) FindByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User
	row := s.db.QueryRowContext(ctx, "SELECT id, login, password, created_at FROM users WHERE login = $1", login)

	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		s.logger.Errorf("failed to find user: %v", err)
		return model.User{}, err
	}

	return user, nil
}
