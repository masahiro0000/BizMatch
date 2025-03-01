package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/masahiro0000/BizMatch/internal/domain"
)

type likeRepositoryImpl struct {
	db *sqlx.DB
}

func NewLikeRepositoryImpl(db *sqlx.DB) domain.LikeRepository {
	return &likeRepositoryImpl{
		db: db,
	}
}

// CreateLike inserts a new like record into the database.
func (r *likeRepositoryImpl) CreateLike(fromUserID, toUserID int64, status string) error{
	now := time.Now()
	query := `
		INSERT INTO likes (from_user_id, to_user_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, fromUserID, toUserID, status, now, now)
	return err
}

// GetLike retrieves a like record into the database.
func (r *likeRepositoryImpl) GetLike(fromUserID, toUserID int64) (*domain.Like, error) {
	var l domain.Like
	query := `
		SELECT * FROM likes
		WHERE from_user_id = $1 AND to_user_id = $2
		LIMIT 1
	`
	err := r.db.Get(&l, query, fromUserID, toUserID)
	if err != nil {
		// Return nil if no record is found.
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}
