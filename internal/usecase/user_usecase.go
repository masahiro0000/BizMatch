package usecase

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/domain"
	"github.com/masahiro0000/BizMatch/internal/infrastructure/repository"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: repo,
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
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, repository.ErrUserNotFound
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