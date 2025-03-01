package domain

type LikeRepository interface {
	CreateLike(fromUserID, toUserID int64, status string) error
	GetLike(fromUserID, toUserID int64) (*Like, error)
}