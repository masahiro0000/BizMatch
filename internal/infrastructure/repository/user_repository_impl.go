package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/masahiro0000/BizMatch/internal/domain"
)

var ErrUserNotFound = errors.New("ユーザーが見つかりません")

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

func (r *userRepositoryImpl) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	query := "SELECT username, password FROM users WHERE username = $1"

	// Execute the SQL query using username and map the result to the user variable.
	if err := r.db.Get(&user, query, username); err != nil {
		// If no rows are returned, return an error indicated the user was not found.
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// If successful, return a pointer to the user.
	return &user, nil
}