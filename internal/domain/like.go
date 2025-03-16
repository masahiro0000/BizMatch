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
	ErrGetLikeFailed = errors.New("いいね情報の取得に失敗しました")
	ErrUpdateLikeFailed = errors.New("いいね情報の更新に失敗しました")
	ErrCreateLikeFailed = errors.New("いいね情報の作成に失敗しました")
	ErrCannotCancelAfterMatch = errors.New("マッチング成立後はいいねをキャンセルできません")
	ErrGetFromUserIDsFailed = errors.New("いいね送信元ユーザーIDの取得に失敗しました")
)