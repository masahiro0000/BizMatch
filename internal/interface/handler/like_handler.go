package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/domain"
	"github.com/masahiro0000/BizMatch/internal/interface/handler/util"
	"github.com/masahiro0000/BizMatch/internal/usecase"
)

type LikeHandler struct {
	likeUsecase *usecase.LikeUsecase
	userUsecase *usecase.UserUsecase
}

func NewLikeHandler(lu *usecase.LikeUsecase, uu *usecase.UserUsecase) *LikeHandler {
	return &LikeHandler{
		likeUsecase: lu,
		userUsecase: uu,
	}
}

// Like processes the like request from a user.
func (h *LikeHandler) Like(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}
	fromUserID := currentUser.ID

	// Parse the target user ID from the URL parameter.
	toUserIDStr := c.Param("id")
	toUserID, err := strconv.ParseInt(toUserIDStr, 10, 64)
	if err != nil {
		log.Printf("fail to parseInt toUserID:%v", err)
		// Render an error message if the user ID cannot be parsed.
		c.HTML(http.StatusInternalServerError, "search_users.html", gin.H{
			"error": domain.ErrUserIDNotFound,
		})
	}

	// Retrieve the user information.
	toUser, err := h.userUsecase.GetUserByID(toUserID)
	if err != nil {
		// Render an error message if user information retrieval fails.
		c.HTML(http.StatusOK, "user_detail.html", gin.H{
			"user": toUser,
			"error": err,
		})
		return
	}

	// Attempt to send a like using the business logic.
	if err := h.likeUsecase.Like(fromUserID, toUserID); err != nil {
		// Render an error message if sending the like fails.
		c.HTML(http.StatusOK, "user_detail.html", gin.H{
			"user": toUser,
			"error": err,
		})
		return
	}

	c.Redirect(http.StatusFound, "/users/detail/"+toUser.Username)
}

// Cancel processes the cancel request from a user.
func (h *LikeHandler) Cancel(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}
	fromUserID := currentUser.ID

	// Parse the target user ID from the URL parameter.
	toUserIDStr := c.Param("id")
	toUserID, err := strconv.ParseInt(toUserIDStr, 10, 64)
	if err != nil {
		log.Printf("fail to parseInt toUserID:%v", err)
		// Render an error message if the user ID cannot be parsed.
		c.HTML(http.StatusInternalServerError, "search_users.html", gin.H{
			"error": domain.ErrUserIDNotFound,
		})
	}

	// Retrieve the user information.
	toUser, err := h.userUsecase.GetUserByID(toUserID)
	if err != nil {
		// Render an error message if user information retrieval fails.
		c.HTML(http.StatusOK, "user_detail.html", gin.H{
			"user": toUser,
			"error": err,
		})
		return
	}

	// Attempt to send a cancel using the business logic.
	err = h.likeUsecase.Cancel(fromUserID, toUserID)
	if err != nil {
		c.HTML(http.StatusOK, "user_detail.html", gin.H{
			"user": toUser,
			"error": err,
		})
		return
	}

	c.Redirect(http.StatusFound, "/users/detail/"+toUser.Username)
}

// ReceivedLikes handles the HTTP request for displaying the list of users who liked the current user.
func (h *LikeHandler) ReceivedLikes(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	// Call the usecase to get the list of users who have liked the current user.
	users, err := h.likeUsecase.GetReceivedLikes(currentUser.ID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "mypage.html", gin.H{
			"user": currentUser,
		})
		return
	}

	c.HTML(http.StatusOK, "receive_likes.html", gin.H{
		"user": currentUser,
		"receive_like_users": users,
	})
}