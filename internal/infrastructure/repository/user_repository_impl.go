package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/masahiro0000/BizMatch/internal/domain"
)

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepositoryImpl(db *sqlx.DB) domain.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

// CreateUser inserts a new user into the database.
func (r *userRepositoryImpl) CreateUser(user *domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, display_name, password) VALUES ($1, $2, $3)",
						user.Username, user.DisplayName, user.HashedPassword)
	if err != nil {
		var pqErr *pq.Error
		// Check if the error is a PostgreSQL error for duplicate entries.
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}
