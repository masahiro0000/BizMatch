package handler

import (
	"encoding/json"
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
	msgUsecase *usecase.MessageUsecase
}

func NewMatchHandler(mu *usecase.MatchUsecase, uu *usecase.UserUsecase, msgu *usecase.MessageUsecase) *MatchHandler {
	return &MatchHandler{
		matchUsecase: mu,
		userUsecase: uu,
		msgUsecase: msgu,
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
	matchIDStr := c.Param("match_id")
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

	// Convert the data to JSON format.
	matchJSON, err := json.Marshal(match)
	if err != nil {
		log.Printf("fail to json:%v", err)
		return
	}
	currentUserJSON, err := json.Marshal(currentUser)
	if err != nil {
		log.Printf("fail to json:%v", err)
		return
	}

	c.HTML(http.StatusOK, "message.html", gin.H{
		"user": currentUser,
		"match": string(matchJSON),
		"userJSON": string(currentUserJSON),
	})
}

func (h *MatchHandler) GetMessages(c *gin.Context) {
	log.Printf("first log")
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}
	fromUserID := currentUser.ID

	// Get the match ID from the parameter.
	matchIDStr := c.Param("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		log.Printf("fail to parse matchID:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid matchID"})
		return
	}

	// Retrieve messages for the specified match ID and user ID.
	messages, err := h.msgUsecase.GetMessages(matchID, fromUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}