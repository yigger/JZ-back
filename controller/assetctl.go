package controller

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/yigger/JZ-back/service"
)

// 资产首页所需数据
func GetWalletsAction(c echo.Context) error {
	return c.JSON(http.StatusOK, service.Asset.GetWallet())
}