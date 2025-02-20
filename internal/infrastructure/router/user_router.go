package router

import (
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/interface/handler"
)

// UserRouter sets up all user-related HTTP routes.
func UserRouter(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	// Route for registration of the user account.
	rg.GET("/signup", userHandler.ShowSignupForm)
	rg.POST("/signup", userHandler.Signup)

	// Route for presenting additional user registration details.
	rg.GET("/register_info", userHandler.ShowRegisterInfo)
	rg.POST("/register_info", userHandler.RegisterInfo)

	// Route for login.
	rg.GET("/login", userHandler.ShowLoginForm)
	rg.POST("/login", userHandler.Login)

	// Route for logout.
	rg.POST("/logout", userHandler.Logout)

	// Route for mypage
	rg.GET("/mypage", userHandler.ShowMypage)
}