package controller

import (
	// "fmt"
	"net/http"
	"github.com/labstack/echo"

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
	return c.JSON(http.StatusOK, service.Statement.GetStatements())
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