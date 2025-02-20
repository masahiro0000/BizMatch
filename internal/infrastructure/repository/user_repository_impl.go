package repository

import (
	"database/sql"
	"errors"
	"log"

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
						user.Username, user.DisplayName, user.Password)
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
	query := "SELECT * FROM users WHERE username = $1"

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

// RegisterInfo updates the user's information in the database.
func (r *userRepositoryImpl) RegisterInfo(user *domain.User) error {
	// SQL query that update the user's information.
	query := `
		UPDATE users
		SET
			username 	 		= :username,
			display_name		= :display_name,
			prefecture_id		= :prefecture_id,
			industry_id 		= :industry_id,
			job_id				= :job_id,
			position_id 		= :position_id,
			age					= :age,
			gender				= :gender,
			photo				= :photo,
			profile_description	= :profile_description
		WHERE
			id = :id
		`

	// Execute the SQL update query.
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		log.Printf("Failed to update user's information. query:%v, error:%v", query, err)
		return err
	}

	return nil
}