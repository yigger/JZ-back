package controllers

import (
	"github.com/labstack/echo"
)

const (
	Json404 = map[string]interface{}{"msg": "登录失败",	"status": 401}
)

func Render404(c echo.Context) {
	// return c.String(http.StatusOK, "login test")
}