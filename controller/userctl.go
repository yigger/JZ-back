package controller

import (
	"github.com/labstack/echo"
	"net/http"

	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/service"
)

// 登录的入口
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

// 获取用户信息
func GetUserAction(c echo.Context) error {
	currentUser := service.CurrentUser
	json := map[string]interface{}{
		"id": currentUser.ID,
		"avatar_url": currentUser.AvatarPath(),
		"nickname": currentUser.Nickname,
		"persist": currentUser.PersistDay(),
		"bonus_points": currentUser.BonusPoints,
		"sts_count": currentUser.StatementCount(),
		"email": currentUser.Email,
		"remind": currentUser.Remind == 1,
		"bg_avatar_url": currentUser.BgAvatarPath(),
		"hidden_asset_money": currentUser.HiddenAssetMoney,
		"position_1_human_name": model.POSITION_ONE_HUMAN[currentUser.HeaderPosition1],
		"position_2_human_name": model.POSITION_TWO_HUMAN[currentUser.HeaderPosition2],
		"position_3_human_name": model.POSITION_THREE_HUMAN[currentUser.HeaderPosition3],
		"position_1_amount": currentUser.GetHeaderAmount(currentUser.HeaderPosition1),
		"position_2_amount": currentUser.GetHeaderAmount(currentUser.HeaderPosition2),
		"position_3_amount": currentUser.GetHeaderAmount(currentUser.HeaderPosition3),
	}
	return c.JSON(http.StatusOK, json)
}

// 更新用户信息
func updateUserAction(c echo.Context) error {
	json := RenderJson()
	defer c.JSON(http.StatusOK, json)

	params := make(map[string]interface{})
	if err := c.Bind(&params); err != nil {
		json.Status = CodeErr
		json.Msg = "err params"
		logger.Info(err)
	}

	userParams := params["user"].(map[string]interface{})
	if _, ok := userParams["alreadyLogin"]; ok {
		userParams["alreadyLogin"] = params["already_login"]
	}
	service.User.UpdateUser(userParams)

	return nil
}