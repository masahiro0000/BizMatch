package repository

import (
	"database/sql"
	"log"
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
func (r *likeRepositoryImpl) CreateLikeRecord(fromUserID, toUserID int64, status string) error{
	now := time.Now()
	query := `
		INSERT INTO likes (from_user_id, to_user_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, fromUserID, toUserID, status, now, now)
	if err != nil {
		log.Printf("fail to create like record: %v", err)
		return err
	}
	return nil
}

// UpdateLikeRecord updates a record into the database.
func (r *likeRepositoryImpl) UpdateLikeRecord(fromUserID, toUserID int64, status string) error {
	now := time.Now()
	query := `
		UPDATE likes
		SET status = $1, updated_at = $2
		WHERE from_user_id = $3 AND to_user_id = $4
	`
	_, err := r.db.Exec(query, status, now, fromUserID, toUserID)
	if err != nil {
		log.Printf("fail to update like record: %v", err)
		return err
	}
	return nil
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
			log.Printf("like record not found")
			return nil, nil
		}
		log.Printf("fail to get like record: %v", err)
		return nil, err
	}
	return &l, nil
}

// GetFromUserIDsWhoLiked receives the IDs of users who have liked the specified user(toUserID).
func (r *likeRepositoryImpl) GetFromUserIDsWhoLiked(toUserID int64) ([]int64, error) {
	// SQL query to select the 'from_user_id' of likes WHERE:
	// - The target user is the one being liked
	// - The like status is 'LIKE'
	// - There is no existing match between the liking user and the target user
	query := `
		SELECT l.from_user_id
		FROM likes l
		WHERE l.to_user_id = $1
			AND status = 'LIKE'
			AND NOT EXISTS (
				SELECT 1
				FROM matches m
				WHERE (
					(m.user1_id = l.from_user_id AND m.user2_id = l.to_user_id)
					OR
					(m.user1_id = l.to_user_id AND m.user2_id = l.from_user_id)
				)
			)
	`
	var fromUserIDs []int64
	if err := r.db.Select(&fromUserIDs, query, toUserID); err != nil {
		log.Printf("fail to get from user IDs who liked: %v", err)
		return nil, err
	}
	return fromUserIDs, nil
}

// GetToUserIDsWhoLiked receives the IDs of users who have been liked by the specified user(fromUserID).
func (r *likeRepositoryImpl) GetToUserIDsWhoLiked(fromUserID int64) ([]int64, error) {
	// SQL query to select the 'to_user_id' of likes WHERE:
	// - The liking user is the one liking
	// - The like status is 'LIKE'
	query := `
		SELECT l.to_user_id
		FROM likes l
		WHERE l.from_user_id = $1
			AND status = 'LIKE'
	`
	var toUserIDs []int64
	if err := r.db.Select(&toUserIDs, query, fromUserID); err != nil {
		return nil, err
	}
	return toUserIDs, nil
}