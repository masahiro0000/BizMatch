package util

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/domain"
)

// SetCurrentUser stores the provided user in the current session.
func SetCurrentUser(c *gin.Context, user *domain.User) () {
	if user != nil {
		// Create a new session and store the authenticated user data.
		session := sessions.Default(c)
		session.Set("user", user)

		// Save the session.
		if err := session.Save(); err != nil {
			c.String(http.StatusInternalServerError, "セッション保存エラーが発生しました")
		}
	}
}

// GetCurrentUser retrieves the current user from the session.
func GetCurrentUser(c *gin.Context) (*domain.User, error) {
	session := sessions.Default(c)
	user := session.Get("user")
	currentUser, ok := user.(*domain.User)

	// If the assertion fails or current user is nil, return an error.
	if !ok || currentUser == nil {
		return nil, errors.New("no valid user session found")
	}

	return currentUser, nil
}
