package controller

import (
	"net/http"
	"github.com/labstack/echo"
)

func GetWalletsAction(c echo.Context) error {
	
	return c.JSON(http.StatusOK, nil)
}