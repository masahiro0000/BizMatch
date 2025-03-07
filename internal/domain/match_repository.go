package domain

type MatchRepository interface {
	CreateMatch(user1ID, user2ID int64, status string) error
	GetMatch(user1ID, user2ID int64) (*Match, error)
	GetMatchByID(matchID int64) (*Match, error)
	GetMatchByUserID(userID int64) ([]*Match, error)
}