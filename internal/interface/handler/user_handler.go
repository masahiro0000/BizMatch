package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/domain"
	"github.com/masahiro0000/BizMatch/internal/infrastructure/repository"
	"github.com/masahiro0000/BizMatch/internal/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler (u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: u,
	}
}

func (h *UserHandler) ShowSignupForm(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{
	})
}

// Signup processes the signup form submission.
func (h *UserHandler) Signup(c *gin.Context) {
	username := c.PostForm("username")
	displayName := c.PostForm("displayName")
	password := c.PostForm("password")

	err := h.userUsecase.Signup(c, username, displayName, password)
	if err != nil {
		var status int
		// Determine the HTTP status code based on the error type.
		switch {
		case errors.Is(err, domain.ErrInvalidInput) || errors.Is(err, domain.ErrPasswordTooShort):
			status = http.StatusBadRequest
		case errors.Is(err, domain.ErrUserAlreadyExists):
			status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
		}

		// Re-render the signup page with an error message and pre-filled form values.
		c.HTML(status, "signup.html", gin.H{
			"error": err.Error(),
			"username": username,
			"displayName": displayName,
		})
		return
	}

	// On successful signup, redirect the user to the registration info page.
	c.Redirect(http.StatusFound, "/users/register_user_info")
}

func (h *UserHandler) ShowRegisterUserInfo(c *gin.Context) {
	c.HTML(http.StatusOK, "register_user_info.html", gin.H{
	})
}

func (h *UserHandler) ShowLoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Attempt to login using username and password.
	_, err := h.userUsecase.Login(c, username, password)
	if err != nil {
		var status int
		// Adjust the HTTP status code according to the type of error.
		switch {
		case errors.Is(err, repository.ErrUserNotFound) || errors.Is(err, domain.ErrIncorrectPassword):
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}

		c.HTML(status, "login.html", gin.H{
			"error": err.Error(),
			"username": username,
		})
	}

	// If successful, redirect the user to their mypage.
	c.Redirect(http.StatusFound, "/users/mypage")
}

func (h *UserHandler) ShowMypage(c *gin.Context) {
	c.HTML(http.StatusOK, "mypage.html", gin.H{
	})
}