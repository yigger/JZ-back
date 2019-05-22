package model

import (
	"time"
	"fmt"
	"github.com/yigger/JZ-back/conf"
)

type User struct {
	Id       				uint64	`json:"id,primary_key"`
	Openid					string
	Email 					string
	// nickname
	// language
	// city
	// province
	// avatar_url
	// country
	SessionKey				string	`form:"session_key" json:"session_key"`
	// gender
	// uid
	// third_session
	// phone
	// budget
	// bg_avatar_url
	// bonus_points
	// header_position_1
	// header_position_2
	// header_position_3
	// bg_avatar_id
	// remind
	// hidden_asset_money
	// already_login
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

func (*User) GetUserByThirdSession(session string) (*User) {
	user := User{}
	if err := db.Where("third_session = ?", session).First(&user).Error; err != nil {
		return nil
	}

	return &user
}

func (User) GetUserByOpenId(openid string) (*User) {
	user := User{}
	if err := db.Where("openid = ?", openid).First(&user).Error; err != nil {
		return nil
	}

	return &user
}

func (User) IsLogin(session string) bool {
	// development do not check login
	if conf.Development() {
		return true
	}

	var User User
	u := User.GetUserByThirdSession(session)
	if u == nil {
		return false
	}

	return u.CacheSessionVal() != ""
}

func (user User) sessionKey() string {
	var key string
	key = fmt.Sprintf("@user_%d_session_key@", user.Id)
	return key
}

// 缓存中的 session 值
func (user User) CacheSessionVal() (string) {
	cacheVal, err := redisCli.Get(user.sessionKey()).Result()
	if err != nil {
		return ""
	}

	return cacheVal
}

func (User) CreateUser(user *User) {
	db.Create(user)
}

func (User) UpdateUser(user *User) {
	db.Save(user)
}