package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/domain"
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