package controllers

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/yigger/JZ-back/utils"
	"github.com/yigger/JZ-back/model"
)

func Login(c echo.Context) error {
	req := c.Request()
	code := req.Header.Get("X-WX-Code")
	res, err := utils.Code2Session(code)
	if err != nil {
		return c.JSON(http.StatusOK, Json404)
	}

	var User model.User
	user := User.GetUserByOpenId(res.OpenID)
	if user == nil {
		// 创建新的用户
		user = User.CreateUser()
	} else {
		// 重新生成 session_key 并存储到缓存中
		
	}

	if session := user.CacheSessionVal(); user != nil && session != "" {
		return c.JSON(http.StatusOK, map[string]interface{}{"session": session,	"status": 200})
	}

	return c.String(http.StatusOK, "login test")
}
