package usecase

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/domain"
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