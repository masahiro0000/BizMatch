package router

import (
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/interface/handler"
)

func LikeRouter(rg *gin.RouterGroup, likeHandler *handler.LikeHandler) {
	// Route for sending like.
	rg.POST("/:id/like", likeHandler.SendLike)
}