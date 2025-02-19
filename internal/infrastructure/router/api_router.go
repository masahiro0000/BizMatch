package router

import (
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/interface/handler"
)

func ApiRouter(rg *gin.RouterGroup, apiHandler *handler.ApiHandler) {
	rg.GET("/prefectures", apiHandler.GetPrefectures)
	rg.GET("/industries", apiHandler.GetIndustries)
	rg.GET("/jobs", apiHandler.GetJobs)
	rg.GET("/positions", apiHandler.GetPositions)
}
