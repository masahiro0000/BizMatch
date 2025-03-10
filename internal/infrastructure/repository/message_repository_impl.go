package repository

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/masahiro0000/BizMatch/internal/domain"
)

type messageRepositoryImpl struct {
	db *sqlx.DB
}

func NewMessageRepositoryImpl(db *sqlx.DB) domain.MessageRepository {
	return &messageRepositoryImpl{
		db: db,
	}
}

// CreateMessage inserts a new message record into the database.
func (r *messageRepositoryImpl) CreateMessage(matchID, fromUserID int64, content string) error {
	query := `
		INSERT INTO messages (match_id, from_user_id, content, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, matchID, fromUserID, content, time.Now())
	if err != nil {
		log.Printf("createMessage")
		return err
	}
	return nil
}

// GetMessagesByMatchID retrieves messages with the match ID.
func (r *messageRepositoryImpl) GetMessagesByMatchID(matchID int64) ([]domain.Message, error) {
	var messages []domain.Message

	query := `
		SELECT * FROM messages
		WHERE match_id = $1
		ORDER BY created_at ASC
	`
	if err := r.db.Select(&messages, query, matchID); err != nil {
		return nil, err
	}
	return messages, nil
}