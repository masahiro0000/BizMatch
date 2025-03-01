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

// SendLike processes the "send like" request from a user.
func (h *LikeHandler) SendLike(c *gin.Context) {
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
		log.Printf("fail to get user by ID:%v", err)
		// Render an error message if user information retrieval fails.
		c.HTML(http.StatusOK, "user_detail.html", gin.H{
			"user": toUser,
			"error": domain.ErrGetUserInfoFailed,
		})
		return
	}

	// Attempt to send a like using the business logic.
	if err := h.likeUsecase.SendLike(fromUserID, toUserID); err != nil {
		log.Printf("fail to send like: %v", err)
		// Render an error message if sending the like fails.
		c.HTML(http.StatusOK, "user_detail.html", gin.H{
			"user": toUser,
			"error": domain.ErrSendLikeFailed,
		})
		return
	}

	c.Redirect(http.StatusFound, "/users/detail/"+toUser.Username)
}