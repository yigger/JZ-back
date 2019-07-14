package controller

import (
	"github.com/labstack/echo"
	"github.com/yigger/JZ-back/middleware"
	"io/ioutil"
	"net/http"
	"os"
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
	echoServer.GET("check_update", func(c echo.Context) error {
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
	// statement.GET("images", "")

	// 用户相关接口
	user := api.Group("/users")
	user.GET("", GetUserAction)
	user.PUT("/update_user", updateUserAction)

	// 资产页相关接口
	wallet := api.Group("/wallet")
	wallet.GET("", GetWalletsAction)
	wallet.GET("/time_line", GetWalletInfoTimeLineAction)
	wallet.GET("/information", GetWalletInfoAction)
	wallet.GET("/statement_list", GetWalletStatementListAction)

	// 设置页相关接口
	settings := api.Group("/settings")
	settings.GET("", SettingIndexAction)
	settings.GET("/about", AboutMeAction)
	settings.POST("/feedback", FeedBackAction)
	// // 面板背景图片列表
	// settings.GET("/covers", "")
	// // 三个位置下拉数据
	// settings.GET("/positions", "")


	// 预算相关接口
	// budget := api.Group("/budgets")
	// budget.GET("", "")
	// budget.GET("/parent", "")

	// message := api.Group("message")
	// message.GET("", "")

	// 分类相关接口
	category := api.Group("/categories")
	category.POST("", CreateCategoryAction)
	category.GET("/category_list", GetCategoryListAction)
	category.GET("/category_childs", GetCategoryListAction)
	category.GET("/category_statements", GetCategoryStatementsAction)
	category.GET("/parent", GetParentCategoriesAction)
	category.GET("/:id", GetCategoryDetail)

	// icon 列表
	icon := echoServer.Group("/api/icons")
	icon.GET("/categories", func(c echo.Context) error {
		data := getIcons("category")
		return c.JSON(http.StatusOK, data)
	})
	icon.GET("/assets", func(c echo.Context) error {
		data := getIcons("asset")
		return c.JSON(http.StatusOK, data)
	})
}

// FIXME: 这里一纬数组就够了，但是为了兼容目前线上的数据，暂时使用二维
func getIcons(iconType string) (res [][]string) {
	// FIXME: 需要缓存，避免每次的 IO 开销
	root, _ := os.Getwd()
	files, _ := ioutil.ReadDir(root + "/public/images/" + iconType)
	var categories []string
	for _, f := range files {
		categories = append(categories, "/images/" + iconType + "/" + f.Name())
	}
	res = append(res, categories)
	return res
}