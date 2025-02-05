package router

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine{

	// Create a default Gin router with Logger and Recovery middleware.
	r := gin.Default()

	// Serve static CSS files from the ./web/css directory when requested via /css URL path.
	r.Static("/css", "./web/css")

	// Serve static JavaScript files from the ./web/js directory when requested via /js URL path.
	r.Static("/js", "./web/js")

	// Load HTML templates from the specified directory.
	// This allows Gin to render HTML files located under ./web/templates.
	r.LoadHTMLGlob("./web/templates/*")

	return r
}
