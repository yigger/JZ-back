package controller

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/yigger/JZ-back/service"
	"net/http"
)

// 资产首页所需数据
func GetWalletsAction(c echo.Context) error {
	return c.JSON(http.StatusOK, service.Asset.GetWallet())
}

func WalletInformationAction(c echo.Context) error {
	assetId := c.QueryParam("asset_id")
	res, err := service.GetAssetInformation(assetId)
	if err != nil {
		log.Info(err)
		return c.JSON(http.StatusOK, nil)
	} else {
		return c.JSON(http.StatusOK, res)
	}
}