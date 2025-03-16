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

func (u *MessageUsecase) GetMessages(matchID, fromUserID int64) ([]domain.Message, error) {
	// Retrieve the match details using the provided match ID.
	match, err := u.matchRepo.GetMatchByID(matchID)
	if err != nil {
		return nil, domain.ErrGetMatchFailed
	}
	// Check if the match exist.
	if match == nil {
		return nil, domain.ErrMatchNotExist
	}
	// Verify that the requesting user is a participant in the match.
	if match.User1ID != fromUserID && match.User2ID != fromUserID {
		return nil, domain.ErrUserNotRelatedMatch
	}

	// Retrieve all messages associated with the match.
	messages, err := u.messageRepo.GetMessagesByMatchID(matchID)
	if err != nil {
		return nil, domain.ErrGetMessageFailed
	}
	return messages, nil
}