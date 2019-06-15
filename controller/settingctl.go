package controller

import (
	"net/http"
	"github.com/labstack/echo"

	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/service"
	. "github.com/yigger/JZ-back/conf"
)

// “我的”页数据
func SettingIndexAction(c echo.Context) error {
	currentUser := service.CurrentUser
	res := map[string]interface{}{
		"user": map[string]interface{}{
			"name": currentUser.Nickname,
			"avatar_url": currentUser.AvatarPath(),
			"bonus_points": currentUser.BonusPoints,
			"already_login": currentUser.AlreadyLogin == 1,
		},
		"version": Conf.Version,
	}
	return c.JSON(http.StatusOK, res)
}

// 关于洁账的数据
func AboutMeAction(c echo.Context) error {
	res := map[string]interface{}{
		"version": Conf.Version,
		"content": Conf.AboutMe["content"],
		"others":  Conf.AboutMe["others"],
	}
	return c.JSON(http.StatusOK, res)
}

func FeedBackAction(c echo.Context) error {
	res := RenderJson()
	defer c.JSON(http.StatusOK, res)

	feedback := &model.Feedback{}
	if err:= c.Bind(&feedback); err != nil {
		res.Msg = "参数绑定失败"
		return nil
	}
	feedback.UserId = uint64(service.CurrentUser.ID)

	var Feedback model.Feedback
	Feedback.Create(feedback)
	return nil
}
