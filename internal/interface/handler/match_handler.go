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

type MatchHandler struct {
	matchUsecase *usecase.MatchUsecase
	userUsecase *usecase.UserUsecase
}

func NewMatchHandler(mu *usecase.MatchUsecase, uu *usecase.UserUsecase) *MatchHandler {
	return &MatchHandler{
		matchUsecase: mu,
		userUsecase: uu,
	}
}

// ListMatches retrieves the list of users who have matches with the current user.
func (h *MatchHandler) ListMatches(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	userID := currentUser.ID

	// Get the match record for the current user.
	matches, err := h.matchUsecase.ListMatch(userID)
	if err != nil {
		log.Printf("fail to get ListMatch:%v", err)
		c.HTML(http.StatusInternalServerError, "mypage.html", gin.H{
			"error": domain.ErrCannotGetMatchList,
		})
		return
	}

	var matchedUsers []interface{}
	var matchIDs []int64
	for _, match := range matches {
		var otherUserID int64
		if match.User1ID == userID {
			otherUserID = match.User2ID
		} else {
			otherUserID = match.User1ID
		}
		// Retrieve user detail for the matched user.
		user, err := h.userUsecase.GetUserByID(otherUserID)
		if err != nil {
			// Skip if unable to retrieve user detail.
			log.Printf("fail to get user by ID:%v", err)
			continue
		}
		matchedUsers = append(matchedUsers, user)
		matchIDs = append(matchIDs, match.ID)
	}

	c.HTML(http.StatusOK, "match_list.html", gin.H{
		"matchedUsers": matchedUsers,
		"matchIDs": matchIDs,
	})
}

func (h *MatchHandler) ShowMatchMessage(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	// Parse the target match ID from the URL parameter.
	matchIDStr := c.Param("matchID")
	matchID, _ := strconv.ParseInt(matchIDStr, 10, 64)

	// Retrieve the match with the match ID.
	match, err := h.matchUsecase.GetMatchByID(matchID)
	if err != nil {
		log.Printf("fail to get match by ID:%v", err)
		c.HTML(http.StatusBadRequest, "mypage.html", gin.H{
			"user": currentUser,
		})
		return
	}

	c.HTML(http.StatusOK, "message.html", gin.H{
		"user": currentUser,
		"match": match,
	})
}