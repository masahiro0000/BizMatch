package router

import (
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/usecase"
	"github.com/masahiro0000/BizMatch/internal/websocket"
)

func WebSocketRouter(r *gin.Engine, hub *websocket.Hub, msgUsecase *usecase.MessageUsecase) {
	// Create a new WebSocket handler with the hub and a callback function.
	wsHandler := websocket.NewWebSocketHandler(
		hub,
		func (matchID, fromUserID int64, content string) error {
			return msgUsecase.SendMessage(matchID, fromUserID, content)
		},
	)

	// Route for registering of WebSocket connections.
	r.GET("/ws/:match_id", wsHandler.ServeWs)
}