package controller

import (
	// "fmt"
	"net/http"
	"github.com/labstack/echo"

	"github.com/yigger/JZ-back/service"
	// "github.com/yigger/JZ-back/logs"
)

func ShowStatementsAction(c echo.Context) error {
	json := RenderJson()
	defer c.JSON(http.StatusOK, json)

	json.Data = service.Statement.GetStatements()

	return nil
}