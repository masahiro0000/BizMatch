package usecase

import (
	"github.com/masahiro0000/BizMatch/internal/domain"
)

type LikeUsecase struct {
	likeRepo	domain.LikeRepository
	matchRepo	domain.MatchRepository
	userRepo	domain.UserRepository
}

func NewLikeUsecase(
	likeRepo domain.LikeRepository,
	matchRepo domain.MatchRepository,
	userRepo domain.UserRepository,
) *LikeUsecase {
	return &LikeUsecase{
		likeRepo:	likeRepo,
		matchRepo:	matchRepo,
		userRepo: 	userRepo,
	}
}

// SendLike processes a like action from one user to another.
func (u *LikeUsecase) Like(fromUserID, toUserID int64) error {
	// Prevent a user from linking themselves.
	if fromUserID == toUserID {
		return domain.ErrSelfLikeAndCancel
	}

	// Check if the like already exists.
	existingLike, err := u.likeRepo.GetLike(fromUserID, toUserID)
	if err != nil {
		return domain.ErrGetLikeFailed
	}
	// Existing like record found
	if existingLike != nil {
		// If the status is "LIKE", return an error the user has already liked target.
		if existingLike.Status == "LIKE" {
			return domain.ErrAlreadyLike
		// If the status is "CANCEL", update the record to "LIKE" to re-enable like.
		} else if existingLike.Status == "CANCEL" {
			if err := u.likeRepo.UpdateLikeRecord(fromUserID, toUserID, "LIKE"); err != nil {
				return domain.ErrUpdateLikeFailed
			}
		}
	// No existing like record found
	} else {
		if err := u.likeRepo.CreateLikeRecord(fromUserID, toUserID, "LIKE"); err != nil {
			return domain.ErrCreateLikeFailed
		}
	}

	// Check if the target user has already liked the current user.
	oppositeLike, err := u.likeRepo.GetLike(toUserID, fromUserID)
	if err != nil {
		return domain.ErrGetLikeFailed
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
			return domain.ErrCannotGetMatchList
		}
		if existingMatch != nil && existingMatch.Status == "ACTIVE" {
			return domain.ErrAlreadyMatch
		}

		// Create a new match since both users have liked each other.
		err = u.matchRepo.CreateMatch(user1ID, user2ID, "ACTIVE")
		if err != nil {
			return domain.ErrCreateMatchFailed
		}
	}
	return nil
}

func (u *LikeUsecase) Cancel(fromUserID, toUserID int64) error {
	// Prevent a user from canceling themselves.
	if fromUserID == toUserID {
		return domain.ErrSelfLikeAndCancel
	}

	// After the match has been created, cannot change like status to "cancel".
	user1ID, user2ID := fromUserID, toUserID
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}
	existingMatch, err := u.matchRepo.GetMatch(user1ID, user2ID)
	if err != nil {
		return domain.ErrGetMatchFailed
	}
	if existingMatch != nil {
		return domain.ErrCannotCancelAfterMatch
	}

	// Returning the existing like record between the two users from the repository.
	existingLike, err  := u.likeRepo.GetLike(fromUserID, toUserID)
	if err != nil {
		return domain.ErrGetLikeFailed
	}
	if existingLike != nil {
		// If the like record already has a "CANCEL" status, it means the like is already canceled.
		if existingLike.Status == "CANCEL" {
			return domain.ErrAlreadyCancel
		// If the record is currently a "LIKE", update it to "CANCEL" to reflect the cancellation.
		} else if existingLike.Status == "LIKE" {
			if err := u.likeRepo.UpdateLikeRecord(fromUserID, toUserID, "CANCEL"); err != nil {
				return domain.ErrUpdateLikeFailed
			}
		}
	} else {
		// If no like record exists, create a new record with a "CANCEL" status.
		if err := u.likeRepo.CreateLikeRecord(fromUserID, toUserID, "CANCEL"); err != nil {
			return domain.ErrCreateLikeFailed
		}
	}
	return nil
}

// GetReceivedLikes retrieves a list of users who have liked the specified user.
func (u *LikeUsecase) GetReceivedLikes(toUserID int64) ([]*domain.User, error) {
	// Retrieve the list of user IDs who have liked the target user.
	fromUserIDs, err := u.likeRepo.GetFromUserIDsWhoLiked(toUserID)
	if err != nil {
		return nil, domain.ErrGetFromUserIDsFailed
	}

	var users []*domain.User
	// For each user ID, fetch the complete user information.
	for _, fromUserID := range fromUserIDs {
		user, err := u.userRepo.GetUserByID(fromUserID)
		if err != nil {
			return nil, domain.ErrGetUserFailed
		}
		users = append(users, user)
	}
	return users, nil
}