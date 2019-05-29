package controller

import (
	// "fmt"
	"net/http"
	"github.com/labstack/echo"

	"github.com/yigger/JZ-back/service"
	// "github.com/yigger/JZ-back/logs"
)

func ShowStatementsAction(c echo.Context) error {
	return c.JSON(http.StatusOK, service.Statement.GetStatements())
}