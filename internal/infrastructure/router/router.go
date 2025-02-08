package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Router encapsulates the Gin engine and the HTTP server
type Router struct {
	engine 	   *gin.Engine
	httpServer *http.Server
}

// NewRouter creates and configures a new router instance
func NewRouter() *Router {
	r := gin.Default()

	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")

	r.LoadHTMLGlob("./web/templates/*")

	r.GET("/top", func(c *gin.Context) {
		c.HTML(http.StatusOK, "top.html", gin.H{
		})
	})

	return &Router{
		engine: r,
	}
}

// Run starts the HTTP server
func (r *Router) Run(addr string) error {
	r.httpServer = &http.Server{
		Addr:		addr,
		Handler:	r.engine,
	}

	return r.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (r *Router) Shutdown(ctx context.Context) error {
	if r.httpServer != nil {
		return r.httpServer.Shutdown(ctx)
	}
	return nil
}