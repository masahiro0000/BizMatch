package domain

import (
	"errors"
	"time"
)

type Like struct {
	ID			int64 `db:"id"`
	FromUserID	int64 `db:"from_user_id"`
	ToUserID	int64 `db:"to_user_id"`
	Status		string `db:"status"`
	CreatedAt	time.Time `db:"created_at"`
	UpdatedAt	time.Time `db:"updated_at"`
}

// Predefined errors related to like functionality.
var (
	ErrSelfLikeAndCancel = errors.New("自分自身にはいいねやキャンセルはできません")
	ErrAlreadyLike = errors.New("既にいいね済みです")
	ErrAlreadyCancel = errors.New("既にキャンセル済みです")
	ErrUserIDNotFound = errors.New("いいね対象のユーザーIDが取得できませんでした。ユーザー検索からやり直してください。")
	ErrGetUserInfoFailed = errors.New("ユーザー情報の取得に失敗しました。ユーザー検索からやり直してください。")
	ErrSendLikeFailed = errors.New("いいね送信に失敗しました")
	ErrCancelLikeNotFound = errors.New("キャンセル対象のいいねが存在しないか、既にキャンセルされています")
	ErrCannotCancelAfterMatch = errors.New("マッチング成立後はいいねをキャンセルできません")
)