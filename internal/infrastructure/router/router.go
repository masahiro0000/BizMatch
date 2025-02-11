package router

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/interface/handler"
)

// Router encapsulates the Gin engine and the HTTP server
type Router struct {
	engine 	   *gin.Engine
	httpServer *http.Server
}

// NewRouter creates and configures a new router instance
func NewRouter(userHandler *handler.UserHandler) *Router {
	r := gin.Default()

	sessionSecretKey := os.Getenv("SESSION_SECRET_KEY")

	// Create a new session using the secret key.
	store := cookie.NewStore([]byte(sessionSecretKey))
	r.Use(sessions.Sessions("bizmatchSession", store))

	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")

	r.LoadHTMLGlob("./web/templates/*")

	r.GET("/top", func(c *gin.Context) {
		c.HTML(http.StatusOK, "top.html", gin.H{
		})
	})

	// Initialize user-related endpoint.
	userGroup := r.Group("/users")
    UserRouter(userGroup, userHandler)

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