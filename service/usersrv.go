package service

import (
	"sync"
	"strconv"

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
	gender, err := strconv.ParseUint(userParams["gender"].(string), 10, 64)
	if nil != err {
		return nil, err
	}

	CurrentUser.Country = userParams["country"].(string)
	CurrentUser.City = userParams["city"].(string)
	CurrentUser.Gender = gender
	CurrentUser.Language = userParams["language"].(string)
	CurrentUser.Province = userParams["province"].(string)
	CurrentUser.Nickname = userParams["nickName"].(string)

	var User model.User
	User.UpdateUser(CurrentUser)
	
	return CurrentUser, nil
}

