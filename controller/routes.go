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
	// 中间件，身份校验
	api.Use(middleware.CheckOpenId)

	// 获取首页的账单列表
	api.GET("/index", ShowStatementsAction)

	// 获取首页的头部信息
	api.GET("/header", ShowIndexHeader)

	// 相关账单
	statement := api.Group("/statements")
	statement.POST("/create", CreateStatementAction)


	// 更新用户
	user := api.Group("/users")
	user.PUT("/update_user", updateUserAction)
}
