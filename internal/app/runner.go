package app

import (
	"context"
	"encoding/gob"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/masahiro0000/BizMatch/internal/domain"
	"github.com/masahiro0000/BizMatch/internal/infrastructure/db"
	"github.com/masahiro0000/BizMatch/internal/infrastructure/repository"
	"github.com/masahiro0000/BizMatch/internal/infrastructure/router"
	"github.com/masahiro0000/BizMatch/internal/interface/handler"
	"github.com/masahiro0000/BizMatch/internal/usecase"
)

func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load environment variables from the .env file.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database connection
	dbConn, err := db.InitDB()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	// Register the User struct with the gob package to enable session storage.
	gob.Register(&domain.User{})

	// Initialize user-related functionality
	userRepo := repository.NewUserRepositoryImpl(dbConn)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	// Initialize api functionality
	apiRepo := repository.NewApiRepositoryImpl(dbConn)
	apiUsecase := usecase.NewApiUsecase(apiRepo)
	apiHandler := handler.NewApiHandler(apiUsecase)

	// Initialize match functionality
	matchRepo := repository.NewMatchRepositoryImpl(dbConn)
	matchUsecase := usecase.NewMatchUsecase(matchRepo)
	matchHandler := handler.NewMatchHandler(matchUsecase, userUsecase)

	// Initialize like functionality
	likeRepo := repository.NewLikeRepositoryImpl(dbConn)
	likeUsecase := usecase.NewLikeUsecase(likeRepo, matchRepo, userRepo)
	likeHandler := handler.NewLikeHandler(likeUsecase, userUsecase)

	// Initialize message functionality
	messageRepo := repository.NewMessageRepositoryImpl(dbConn)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, matchRepo)

	// Create a new router instance
	r := router.NewRouter(userHandler, apiHandler, likeHandler, matchHandler, messageUsecase)

	errCh := make(chan error, 1)
	// Start the HTTP server in a separate goroutine
	go func() {
		errCh <- r.Run(":8080")
	}()

	// Wait for either a shutdown signal or an error from HTTP server
	select {
		case <- ctx.Done():
			return r.Shutdown(ctx)
		case err := <-errCh:
			return err
	}
}