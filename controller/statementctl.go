package controller

import (
	"net/http"
	"github.com/labstack/echo"
	
	"github.com/yigger/JZ-back/utils"
	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/service"
)

func ShowIndexHeader(c echo.Context) error {
	currentUser := service.CurrentUser
	var json = map[string]interface{}{
		"bg_avatar": currentUser.BgAvatarPath(),
		"has_no_read": service.CurrentUser.WaitReadMessage(),
		"show_notice_bar": service.CurrentUser.ShowNoticeBar(),
		"notice_bar_path": nil,
		"notice_text": nil,
		"position_1_human_name": model.POSITION_ONE_HUMAN[currentUser.HeaderPosition1],
		"position_2_human_name": model.POSITION_TWO_HUMAN[currentUser.HeaderPosition2],
		"position_3_human_name": model.POSITION_THREE_HUMAN[currentUser.HeaderPosition3],
		"position_1_amount": utils.FormatMoney(currentUser.GetHeaderAmount(currentUser.HeaderPosition1)),
		"position_2_amount": utils.FormatMoney(currentUser.GetHeaderAmount(currentUser.HeaderPosition2)),
		"position_3_amount": utils.FormatMoney(currentUser.GetHeaderAmount(currentUser.HeaderPosition3)),
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

func UpdateStatementAction(c echo.Context) error {
	var result interface{}
	defer c.JSON(http.StatusOK, result)

	statementId := c.Param("id")
	db := model.ConnectDB()
	var statement model.Statement
	if err := db.Model(service.CurrentUser).Where("statements.id = ?", statementId).Association("Statements").Find(&statement).Error; err != nil {
		 result = err
	} else {
		params := make(map[string]interface{})
		if err := c.Bind(&params); err != nil {
			result = "无效的参数"
			return nil
		}

		service.Statement.UpdateStatement(&statement, params)
		result = "xixi"	
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