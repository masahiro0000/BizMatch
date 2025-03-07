package router

import (
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/interface/handler"
)

func MatchRouter(rg *gin.RouterGroup, matchHandler *handler.MatchHandler) {
	// Route for matching list.
	rg.GET("/list", matchHandler.ListMatches)

	// Route for message display.
	rg.GET("/:matchID/message", matchHandler.ShowMatchMessage)
}