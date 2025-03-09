package domain

type MessageRepository interface {
	CreateMessage(matchID, fromUserID int64, content string) error
}