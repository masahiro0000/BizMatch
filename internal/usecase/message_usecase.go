package usecase

import (
	"github.com/masahiro0000/BizMatch/internal/domain"
)

type MessageUsecase struct {
	messageRepo	domain.MessageRepository
	matchRepo	domain.MatchRepository
}

func NewMessageUsecase(msgRepo domain.MessageRepository, matchRepo domain.MatchRepository) *MessageUsecase {
	return &MessageUsecase{
		messageRepo: msgRepo,
		matchRepo: matchRepo,
	}
}

func (u *MessageUsecase) SendMessage(matchID, fromUserID int64, content string) error {
	// Check if the content is empty.
	if content == "" {
		return domain.ErrEmptyMessage
	}

	// Retrieve the match by its ID.
	match, err := u.matchRepo.GetMatchByID(matchID)
	if err != nil {
		return err
	}
	// Check if the match is not found.
	if match == nil {
		return domain.ErrMatchNotExist
	}
	// Verify that the sender is one of the participants of the match.
	if match.User1ID != fromUserID && match.User2ID != fromUserID {
		return domain.ErrUserNotRelatedMatch
	}

	// Create the message using the message repository.
	if err := u.messageRepo.CreateMessage(matchID, fromUserID, content); err != nil {
		return err
	}
	return nil
}