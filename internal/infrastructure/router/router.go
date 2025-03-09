package router

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/interface/handler"
	"github.com/masahiro0000/BizMatch/internal/usecase"
	"github.com/masahiro0000/BizMatch/internal/websocket"
)

// Router encapsulates the Gin engine and the HTTP server
type Router struct {
	engine 	   *gin.Engine
	httpServer *http.Server
}

// NewRouter creates and configures a new router instance
func NewRouter(
	userHandler *handler.UserHandler, apiHandler *handler.ApiHandler, likeHandler *handler.LikeHandler,
	matchHandler *handler.MatchHandler, messageUsecase *usecase.MessageUsecase,
	) *Router {
	r := gin.Default()

	sessionSecretKey := os.Getenv("SESSION_SECRET_KEY")

	// Create a new session using the secret key.
	store := cookie.NewStore([]byte(sessionSecretKey))
	r.Use(sessions.Sessions("bizmatchSession", store))

	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")

	// This configuration is used to store and serve user profile photos and other static assets.
	r.Static("/static", "./web/static")

	r.LoadHTMLGlob("./web/templates/*")

	r.GET("/top", func(c *gin.Context) {
		c.HTML(http.StatusOK, "top.html", gin.H{
		})
	})

	// Initialize user-related endpoint.
	userGroup := r.Group("/users")
    UserRouter(userGroup, userHandler)

	// Routing for like-related function.
	LikeRouter(userGroup, likeHandler)

	// Routing for match-related function.
	matchGroup := r.Group("/match")
	MatchRouter(matchGroup, matchHandler)

	// Routing for sending message.
	hub := websocket.NewHub()
	go hub.Run()
	WebSocketRouter(r, hub, messageUsecase)

	//Initialize api endpoint.
	apiGroup := r.Group("/api")
	ApiRouter(apiGroup, apiHandler)

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