package controller

import (
	"net/http"
	"github.com/labstack/echo"

	"github.com/yigger/JZ-back/service"
)

func Login(c echo.Context) error {
	req := c.Request()
	code := req.Header.Get("X-WX-Code")
	
	json := RenderJson()
	defer c.JSON(http.StatusOK, json)

	if code == "" {
		json.Status = CodeErr
	} else {
		user, err := service.User.Login(code)
		if err != nil {
			json.Msg = "登录失败"
		} else {
			json.Data = user.SessionKey
		}
	}

	return nil
}
