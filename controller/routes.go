package controller

import (
	"net/http"
	"github.com/labstack/echo"

	"github.com/yigger/JZ-back/middleware"
)

var echoServer *echo.Echo

// Start Server
func EchoNew() *echo.Echo {
	echoServer = echo.New()
	loadRoutes()
	echoServer.Static("/", "public")
	return echoServer
}

func loadRoutes() {
	echoServer.GET("check_update", func(c echo.Context) (error) {
		return c.JSON(http.StatusOK, 0)
	})

	api := echoServer.Group("/api")
	// 中间件，身份校验
	api.Use(middleware.CheckOpenId)

	// 登录
	api.POST("/check_openid", LoginAction)

	// 获取首页的账单列表
	api.GET("/index", ShowStatementsAction)

	// 获取首页的头部信息
	api.GET("/header", ShowIndexHeader)

	// 账单相关接口
	statement := api.Group("/statements")
	statement.POST("", CreateStatementAction)
	statement.PUT("/:id", UpdateStatementAction)
	statement.GET("/assets", GetStatementAssetsAction)
	statement.GET("/categories", GetStatementCategoriesAction)
	statement.GET("/category_frequent", CategoryFrequentUseAction)
	statement.GET("/asset_frequent", AssetFrequentUseAction)

	// 用户相关接口
	user := api.Group("/users")
	user.GET("", GetUserAction)
	user.PUT("/update_user", updateUserAction)

	// 资产页相关接口
	wallet := api.Group("/wallet")
	wallet.GET("", GetWalletsAction)

	// 设置页
	settings := api.Group("/settings")
	settings.GET("", SettingIndexAction)
	settings.GET("/about", AboutMeAction)

	category := api.Group("/categories")
	category.GET("/category_list", GetCategoryListAction)
	category.GET("/category_childs", GetCategoryListAction)
	category.GET("/category_statements", GetCategoryStatementsAction)
}
