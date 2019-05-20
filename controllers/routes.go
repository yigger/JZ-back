package controllers

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
	// middleware, check user is login
	echoServer.Use(middleware.CheckOpenId)

	echoServer.POST("/login", Login)

}