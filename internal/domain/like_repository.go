package domain

type LikeRepository interface {
	CreateLikeRecord(fromUserID, toUserID int64, status string) error
	UpdateLikeRecord(fromUserID, toUserID int64, status string) error
	GetLike(fromUserID, toUserID int64) (*Like, error)
	GetFromUserIDsWhoLiked(toUserID int64) ([]int64, error)
	GetToUserIDsWhoLiked(fromUserID int64) ([]int64, error)
}