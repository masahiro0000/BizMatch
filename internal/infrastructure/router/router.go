package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine{
	r := gin.Default()

	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")

	r.LoadHTMLGlob("./web/templates/*")

	r.GET("/top", func(c *gin.Context) {
		c.HTML(http.StatusOK, "top.html", gin.H{
		})
	})

	return r
}
