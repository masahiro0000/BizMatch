package usecase

import (
	"github.com/masahiro0000/BizMatch/internal/domain"
)

type LikeUsecase struct {
	likeRepo	domain.LikeRepository
	matchRepo	domain.MatchRepository
}

func NewLikeUsecase(
	likeRepo domain.LikeRepository,
	matchRepo domain.MatchRepository,
) *LikeUsecase {
	return &LikeUsecase{
		likeRepo:	likeRepo,
		matchRepo:	matchRepo,
	}
}

// SendLike processes a like action from one user to another.
func (u *LikeUsecase) SendLike(fromUserID, toUserID int64) error {
	// Prevent a user from linking themselves.
	if fromUserID == toUserID {
		return domain.ErrSelfLike
	}

	// Check if the like already exists.
	existingLike, err := u.likeRepo.GetLike(fromUserID, toUserID)
	if err != nil {
		return err
	}
	if existingLike != nil && existingLike.Status == "LIKE" {
		return domain.ErrAlreadyLike
	}

	// Create a new like record.
	err = u.likeRepo.CreateLike(fromUserID, toUserID, "LIKE")
	if err != nil {
		return err
	}

	// Check if the target user has already liked the current user.
	oppositeLike, err := u.likeRepo.GetLike(toUserID, fromUserID)
	if err != nil {
		return err
	}
	if oppositeLike != nil && oppositeLike.Status == "LIKE" {
		// Order user IDs to ensure consistency and avoid duplicate matches.
		user1ID, user2ID := fromUserID, toUserID
		if user1ID > user2ID {
			user1ID, user2ID = user2ID, user1ID
		}

		// Check if the match already exists.
		existingMatch, err := u.matchRepo.GetMatch(user1ID, user2ID)
		if err != nil {
			return err
		}
		if existingMatch != nil && existingMatch.Status == "ACTIVE" {
			return domain.ErrAlreadyMatch
		}

		// Create a new match since both users have liked each other.
		err = u.matchRepo.CreateMatch(user1ID, user2ID, "ACTIVE")
		if err != nil {
			return err
		}
	}
	return nil
}