package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/masahiro0000/BizMatch/internal/infrastructure/db"
	"github.com/masahiro0000/BizMatch/internal/infrastructure/repository"
	"github.com/masahiro0000/BizMatch/internal/infrastructure/router"
	"github.com/masahiro0000/BizMatch/internal/interface/handler"
	"github.com/masahiro0000/BizMatch/internal/usecase"
)

func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize the database connection
	dbConn, err := db.InitDB()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	// Initialize user-related functionality
	userRepo := repository.NewUserRepositoryImpl(dbConn)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	// Create a new router instance
	r := router.NewRouter(userHandler)

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