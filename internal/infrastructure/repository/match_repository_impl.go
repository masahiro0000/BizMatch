package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/masahiro0000/BizMatch/internal/domain"
)

type matchRepositoryImpl struct {
	db *sqlx.DB
}

func NewMatchRepositoryImpl(db *sqlx.DB) domain.MatchRepository {
	return &matchRepositoryImpl{
		db: db,
	}
}

// CreateMatch insert a new match record into the database.
func (r *matchRepositoryImpl) CreateMatch(user1ID, user2ID int64, status string) error {
	now := time.Now()
	query := `
		INSERT INTO matches
		(user1_id, user2_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, user1ID, user2ID, status, now, now)
	if err != nil {
		return err
	}
	return nil
}

// GetMatch retrieves a match record from the database based on user IDs.
func (r *matchRepositoryImpl) GetMatch(user1ID, user2ID int64) (*domain.Match, error) {
	var m domain.Match
	query := `
		SELECT * FROM matches
		WHERE user1_id = $1 AND user2_id = $2
		LIMIT 1
	`
	err := r.db.Get(&m, query, user1ID, user2ID)
	if err != nil {
		// Return nil if no record is found.
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}