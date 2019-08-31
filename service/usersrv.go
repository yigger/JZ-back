package service

import (
	"sync"
	// "strconv"

	"github.com/yigger/JZ-back/utils"
	"github.com/yigger/JZ-back/conf"
	"github.com/yigger/JZ-back/model"
)

var CurrentUser = &model.User{}

var User = &userService{
	mutex: &sync.Mutex{},
}

type userService struct {
	mutex *sync.Mutex
}

// Middleware check user login and set global current_user
func (srv *userService) CheckLogin(session string) bool {
	var User model.User

	if conf.Development() {
		CurrentUser = User.GetFirst()	
		return true
	}
	
	CurrentUser = User.GetUserByThirdSession(session)
	if CurrentUser == nil {
		return false
	}

	return CurrentUser.CacheSessionVal() != ""
}

func (srv *userService) Login(code string) (user *model.User, err error) {
	res, err := utils.Code2Session(code)
	if err != nil {
		return
	}

	var User model.User
	user = User.GetUserByOpenId(res.OpenID)
	if user == nil {
		user = &model.User{Openid: res.OpenID, SessionKey: res.SessionKey}
		User.CreateUser(user)
	} else {
		user.SessionKey = res.SessionKey
		User.UpdateUser(user)
	}

	return
}

func (srv *userService) UpdateUser(userParams map[string]interface{}) (*model.User, error) {
	// set user login status
	if _, ok := userParams["alreadyLogin"]; ok && userParams["alreadyLogin"].(bool) {
		CurrentUser.AlreadyLogin = 1
	} else {
		CurrentUser.AlreadyLogin = 0
	}

	if _, ok := userParams["country"]; ok {
		CurrentUser.Country = userParams["country"].(string)
	}

	if _, ok := userParams["city"]; ok {
		CurrentUser.City = userParams["city"].(string)
	}

	if _, ok := userParams["gender"]; ok {
		CurrentUser.Gender = uint64(userParams["gender"].(float64))
	}

	if _, ok := userParams["province"]; ok {
		CurrentUser.Province = userParams["province"].(string)
	}

	if _, ok := userParams["nickName"]; ok {
		CurrentUser.Nickname = userParams["nickName"].(string)
	}

	if _, ok := userParams["avatarUrl"]; ok {
		CurrentUser.AvatarUrl = userParams["avatarUrl"].(string)
	}

	if _, ok := userParams["hidden_asset_money"]; ok {
		CurrentUser.HiddenAssetMoney = uint64(userParams["hidden_asset_money"].(float64))
	}

	var User model.User
	User.UpdateUser(CurrentUser)
	
	return CurrentUser, nil
}

