package domain

import (
	"errors"
	"time"
)

type Match struct {
	ID			int64 `db:"id"`
	User1ID		int64 `db:"user1_id"`
	User2ID		int64 `db:"user2_id"`
	Status		string `db:"status"`
	CreatedAt	time.Time `db:"created_at"`
	UpdatedAt	time.Time `db:"updated_at"`
}

// Predefined errors related to match functionality.
var (
	ErrAlreadyMatch = errors.New("既にマッチしています")
	ErrCannotGetMatchList = errors.New("マッチしているユーザー一覧を取得できませんでした")
	ErrCreateMatchFailed = errors.New("マッチ情報の作成に失敗しました")
	ErrGetMatchFailed = errors.New("マッチ情報の取得に失敗しました")
	ErrUpdateMatchFailed = errors.New("マッチ情報の更新に失敗しました")
)