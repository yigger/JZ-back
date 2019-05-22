package service

import (
	"sync"

	"github.com/yigger/JZ-back/utils"
	"github.com/yigger/JZ-back/model"
	// "fmt"
)

var User = &userService{
	mutex: &sync.Mutex{},
}

type userService struct {
	mutex *sync.Mutex
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


