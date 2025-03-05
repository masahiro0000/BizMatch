package router

import (
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/interface/handler"
)

func MatchRouter(rg *gin.RouterGroup, matchHandler *handler.MatchHandler) {
	rg.GET("/list", matchHandler.ListMatches)
}