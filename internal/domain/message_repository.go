package domain

type MessageRepository interface {
	CreateMessage(matchID, fromUserID int64, content string) error
	GetMessagesByMatchID(matchID int64) ([]Message, error)
}