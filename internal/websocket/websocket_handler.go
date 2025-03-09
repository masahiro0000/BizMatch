package websocket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/masahiro0000/BizMatch/internal/interface/handler/util"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub				*Hub
	handleUsecase	func(matchID, fromUserID int64, content string) error
}

func NewWebSocketHandler(hub *Hub, messageUsecaseFunc func(matchID, fromUserID int64, content string) error) *WebSocketHandler {
	return &WebSocketHandler {
		hub: hub,
		handleUsecase: messageUsecaseFunc,
	}
}

// ServeWs upgrades the HTTP connection to a WebSocket and sets up the client for communication.
func (h *WebSocketHandler) ServeWs(c *gin.Context) {
	// Retrieve the current user from the session.
	user, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	userID := user.ID
	matchIDStr := c.Param("match_id")
	// Convert the match ID string to an int64.
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid matchID",
		})
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		Hub:		h.hub,
		Conn:		conn,
		UserID:		userID,
		MatchID:	matchID,
		send:		make(chan BroadcastMessage),
	}

	// Register the new client with the hub.
	h.hub.register <- client

	// Start a goroutine to continuously read messages from the WebSocket.
	go client.ReadPump(func(c *Client, msg string) error {
		if err := h.handleUsecase(c.MatchID, c.UserID, msg); err != nil {
			return err
		}

		h.hub.Broadcast(BroadcastMessage{
			MatchID:	c.MatchID,
			FromUser:	c.UserID,
			Content:	msg,
		})
		return nil
	})

	// Start a goroutine to continuously write messages to the WebSocket.
	go client.WritePump()
}