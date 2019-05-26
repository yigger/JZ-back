package controller

import (
	"github.com/labstack/echo"

	"github.com/yigger/JZ-back/middleware"
)

var echoServer *echo.Echo

// Start Server
func EchoNew() *echo.Echo {
	echoServer = echo.New()
	loadRoutes()
	return echoServer
}

func loadRoutes() {
	echoServer.POST("/login", LoginAction)

	api := echoServer.Group("/api")
	// middleware, check user is login
	api.Use(middleware.CheckOpenId)

	api.GET("/index", ShowStatementsAction)

	user := api.Group("/users")
	user.PUT("/update_user", updateUserAction)
}
