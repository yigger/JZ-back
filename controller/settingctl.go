package controller

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/yigger/JZ-back/service"
)

func SettingIndexAction(c echo.Context) error {
	currentUser := service.CurrentUser
	res := map[string]interface{}{
		"user": map[string]interface{}{
			"name": currentUser.Nickname,
			"avatar_url": currentUser.AvatarPath(),
			"bonus_points": currentUser.BonusPoints,
			"already_login": currentUser.AlreadyLogin == 1,
		},
		"version": 3.7,
	}
	return c.JSON(http.StatusOK, res)
}