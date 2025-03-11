package router

import (
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/interface/handler"
)

func MatchRouter(rg *gin.RouterGroup, matchHandler *handler.MatchHandler) {
	// Route for matching list.
	rg.GET("/list", matchHandler.ListMatches)

	// Route for message display.
	rg.GET("/:match_id/message", matchHandler.ShowMatchMessage)

	// Route for getting messages.
	rg.GET("/:match_id/messages", matchHandler.GetMessages)

	// Route for block matching.
	rg.POST("/:match_id/block", matchHandler.BlockMatch)
}