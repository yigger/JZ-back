package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/yigger/JZ-back/service"
	"net/http"
)

// 资产首页所需数据
func GetWalletsAction(c echo.Context) error {
	return c.JSON(http.StatusOK, service.Asset.GetWallet())
}

func GetWalletInfoAction(c echo.Context) error {
	assetId := c.QueryParam("asset_id")
	res, err := service.Asset.GetAssetInformation(assetId)
	if err != nil {
		log.Info(err)
		return c.JSON(http.StatusOK, nil)
	} else {
		return c.JSON(http.StatusOK, res)
	}
}

func GetWalletInfoTimeLineAction(c echo.Context) error {
	assetId := c.QueryParam("asset_id")
	data, err := service.Asset.GetAssetTimeLine(assetId)
	if err != nil {
		return c.JSON(http.StatusOK, err)
	} else {
		res := RenderJson()
		res.Data = data
		return c.JSON(http.StatusOK, res)
	}
}

func GetWalletStatementListAction(c echo.Context) error {
	assetId := c.QueryParam("asset_id")
	year := c.QueryParam("year")
	month := c.QueryParam("month")
	data, err := service.Asset.GetStatementsByAsset(assetId, year, month)
	if err != nil {
		return c.JSON(http.StatusOK, err)
	} else {
		res := RenderJson()
		res.Data = data
		return c.JSON(http.StatusOK, res)
	}
}

// 更新账户的余额
func UpdateSurplusAction(c echo.Context) error {
	json := RenderJson()
	defer c.JSON(http.StatusOK, json)
	assetId := c.FormValue("asset_id")
	params := map[string]interface{}{
		"isMess": c.FormValue("is_mess"),
		"amount": c.Param("amount"),
	}
	fmt.Println("haha")
	fmt.Println(c.FormValue("amount"))

	_, err := service.Asset.UpdateSurplus(assetId, params)
	if err != nil {
		json.Msg = err.Error()
	}
	return nil
}