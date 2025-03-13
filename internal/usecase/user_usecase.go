package usecase

import (
	"errors"
	"log"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/domain"
)

type UserUsecase struct {
	userRepo domain.UserRepository
	likeRepo domain.LikeRepository
}

func NewUserUsecase(userRepo domain.UserRepository, likeRepo domain.LikeRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
		likeRepo: likeRepo,
	}
}

// Signup handles the user registration process.
func (u *UserUsecase) Signup(c *gin.Context, username, displayName, password string) error{
	// Create a new user entity.
	user, err := domain.NewUser(username, displayName, password)
	if err != nil {
		// Return the appropriate error based on the type of validation failure.
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			return domain.ErrInvalidInput
		case errors.Is(err, domain.ErrPasswordTooShort):
			return domain.ErrPasswordTooShort
		default:
			return domain.ErrInvalidInput
		}
	}

	// Attempt to store the new user in the repository.
	err = u.userRepo.CreateUser(user)
	if err != nil {
		// If the error indicates that the user already exists, return that error.
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (u *UserUsecase) Login(c *gin.Context, username, password string) (*domain.User, error) {
	var user *domain.User
	var err error

	// Attempt to retrieve the user by username from the repository.
	user, err = u.userRepo.GetUserByUsername(username)
	if err != nil {
		// Return different error message based on the type of error encountered.
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return nil, domain.ErrUserNotFound
		default:
			return nil, errors.New("ログイン中にエラーが発生しました")
		}
	}

	// Verify the provided password against the user's stored credentials.
	if err := user.VerifyPassword(password); err != nil {
		return nil, domain.ErrIncorrectPassword
	}

	return user, nil
}

func (u *UserUsecase) RegisterInfo(user *domain.User) error {
	// Check username and display name.
	if err := domain.CheckUsernameDisplayName(user.Username, user.DisplayName); err != nil {
		log.Printf("Validation failed for username and display name. username:%v ,display name:%v, error:%v",
			user.Username, user.DisplayName, err)
		return err
	}

	// Attempt to update the user's info in the repository.
	if err := u.userRepo.RegisterInfo(user); err != nil {
		log.Printf("Failed to update user info in repository. user ID:%d, error:%v", user.ID, err)
		return err
	}
	return nil
}

func (u *UserUsecase) UpdatePassword(userID int64, oldPassword, newPassword, confirmPassword string) error {
	// Retrieve the user by ID from the repository.
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify that the provided old password matches the user's current password.
	if err := user.VerifyPassword(oldPassword); err != nil {
		return domain.ErrOldPasswordMismatch
	}

	// Check if the new password matches the confirmation password.
	if newPassword != confirmPassword {
		return domain.ErrNewPasswordMismatch
	}

	// Ensure that new password meets the minimum length requirement.
	if len(newPassword) < domain.MinPasswordLength {
		return domain.ErrPasswordTooShort
	}

	// Hash the new password for security.
	hashedPassword, err := domain.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update the user's password in the repository with the hashed new password.
	err = u.userRepo.UpdatePassword(user.ID, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) SearchUsers(filter *domain.UserSearchFilter) ([]*domain.User, error) {
	return u.userRepo.SearchUsers(filter)
}

func (u *UserUsecase) GetUserByUsername(username string) (*domain.User, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil{
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) GetUserByID(userID int64) (*domain.User, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetRecommendedUsers retrieves all users except the user with the provided ID.
func (u *UserUsecase) GetRecommendedUsers(userID int64) ([]*domain.User, error) {
	var candidates []*domain.User
	// Retrieve all users from the repository.
	allUsers, err := u.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// Retrieve the IDs of users who have liked the current user.
	likedUserIDs, err := u.likeRepo.GetToUserIDsWhoLiked(userID)
	if err != nil {
		likedUserIDs = []int64{}
	}

	// Create a map of liked user IDs for faster lookup.
	likedMap := make(map[int64]bool)
	for _, id := range likedUserIDs {
		likedMap[id] = true
	}

	for _, user := range allUsers {
		// Exclude the user with the provided ID from the list of candidates.
		if user.ID == userID {
			continue
		}
		// Exclude users who have already liked the current user.
		if likedMap[user.ID] {
			continue
		}
		candidates = append(candidates, user)
	}
	return candidates, nil
}

// ScoreUsers calculates the score for each candidate user based on the current user's information.
func (u *UserUsecase) ScoreUsers(currentUser *domain.User, candidates []*domain.User) []*domain.ScoredUser {
	var scoredUsers []*domain.ScoredUser

	// Calculate the score for each candidate user.
	for _, candidate := range candidates {
		score := 0
		// Compare the user's information with the candidate's information.
		if currentUser.Prefecture != nil && candidate.Prefecture != nil {
			if *currentUser.Prefecture == *candidate.Prefecture {
				score += 6
			}
		}
		if currentUser.Industry != nil && candidate.Industry != nil {
			if *currentUser.Industry == *candidate.Industry {
				score += 5
			}
		}
		if currentUser.Job != nil && candidate.Job != nil {
			if *currentUser.Job == *candidate.Job {
				score += 4
			}
		}
		if currentUser.Position != nil && candidate.Position != nil {
			if *currentUser.Position == *candidate.Position {
				score += 3
			}
		}
		if currentUser.Age != nil && candidate.Age != nil {
			if GetAgeGroup(*currentUser.Age) == GetAgeGroup(*candidate.Age) {
				score += 2
			}
		}
		if currentUser.Gender != nil && candidate.Gender != nil {
			if *currentUser.Gender == *candidate.Gender {
				score += 1
			}
		}
		scoredUsers = append(scoredUsers, &domain.ScoredUser{
			User: candidate,
			Score: score,
		})
	}

	// Sort the scored users by score in descending order.
	sort.Slice(scoredUsers, func(i, j int) bool {
		return scoredUsers[i].Score > scoredUsers[j].Score
	})

	return scoredUsers
}

func GetAgeGroup(age int64) int64 {
	return age / 10
}