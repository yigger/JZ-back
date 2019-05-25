package controller

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo"

	"github.com/yigger/JZ-back/service"
	// "github.com/yigger/JZ-back/logs"
)

func LoginAction(c echo.Context) error {
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

func updateUserAction(c echo.Context) error {
	json := RenderJson()
	defer c.JSON(http.StatusOK, json)

	params := make(map[string]interface{})
	if err := c.Bind(&params); err != nil {
		json.Status = CodeErr
		json.Msg = "err params"
		fmt.Println(err)
	}
	
	userParams := params["user"].(map[string]interface{})
	service.User.UpdateUser(userParams)
	
	
	// json.Data = user

	return nil
}