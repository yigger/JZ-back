package controller

import (
	"net/http"
	"github.com/labstack/echo"
	
	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/service"
	// "github.com/leekchan/accounting"
	// "github.com/yigger/JZ-back/logs"
)

func ShowIndexHeader(c echo.Context) error {
	// ac := accounting.Accounting{Symbol: "$", Precision: 2}
	var json = map[string]interface{}{
		"bg_avatar": nil,
		"has_no_read": service.CurrentUser.WaitReadMessage(),
		"show_notice_bar": service.CurrentUser.ShowNoticeBar(),
		"notice_bar_path": nil,
		"notice_text": nil,
		"position_1_human_name": nil,
		"position_2_human_name":  nil,
		"position_3_human_name": nil,
		"position_1_amount": nil,
		"position_2_amount": nil,
		"position_3_amount": nil,
	}

	return c.JSON(http.StatusOK, json)
}

func ShowStatementsAction(c echo.Context) error {
	data := service.Statement.GetStatements()
	return c.JSON(http.StatusOK, data)
}

func CreateStatementAction(c echo.Context) error {
	json := RenderJson()
	defer c.JSON(http.StatusOK, json)
	params := make(map[string]interface{})
	if err := c.Bind(&params); err != nil {
		return nil
	}

	statementParams := params["statement"].(map[string]interface{})
	statement, err := service.Statement.CreateStatement(statementParams)
	if err != nil {
		json.Status = 200
		json.Error = err
		return nil
	} else {
		json.Data = statement
	}

	
	return nil
}

func GetStatementAssetsAction(c echo.Context) error {
	data := service.Statement.GetStatementAssets()

	var frequents []map[string]interface{}
	for _, asset := range data["frequent"].([]model.Asset) {
		tmp := map[string]interface{}{
			"id": asset.ID,
			"name": asset.Name,
			"icon_path": asset.IconUrl(),
			"amount": asset.AmountHuman(),
		}
		frequents = append(frequents, tmp)
	}
	data["frequent"] = frequents

	return c.JSON(http.StatusOK, data)
}

func GetStatementCategoriesAction(c echo.Context) error {
	data := service.Statement.GetStatementCategories(c.FormValue("type"))
	
	var frequents []map[string]interface{}
	for _, category := range data["frequent"].([]model.Category) {
		tmp := map[string]interface{}{
			"id": category.ID,
			"name": category.Name,
			"icon_path": category.IconUrl(),
		}
		frequents = append(frequents, tmp)
	}
	data["frequent"] = frequents

	return c.JSON(http.StatusOK, data)
}

func CategoryFrequentUseAction(c echo.Context) error {
	data := service.Statement.CategoryFrequentUse(c.FormValue("type"))
	return c.JSON(http.StatusOK, data)
}

func AssetFrequentUseAction(c echo.Context) error {
	data := service.Statement.AssetFrequentUse()
	return c.JSON(http.StatusOK, data)
}