package domain

import (
	"errors"
	"time"
)

type Message struct {
	ID			int64 `db:"id"`
	MatchID		int64 `db:"match_id"`
	FromUserID	int64 `db:"from_user_id"`
	Content		string `db:"content"`
	CreatedAt	time.Time `db:"created_at"`
	ReadAt		*time.Time `db:"read_at"`
}

var (
	ErrEmptyMessage = errors.New("メッセージが空です")
	ErrMatchNotExist = errors.New("マッチが存在しません")
	ErrUserNotRelatedMatch = errors.New("マッチに関係ないユーザーです")
)