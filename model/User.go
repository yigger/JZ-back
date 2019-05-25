package model

import (
	"fmt"
)

type User struct {
	CommonModel
	Openid					string  `gorm:"unique;not null" json:"open_id,omitempty"`
	Email					string	`gorm:"-"`
	Nickname				string	`json:"nickName,omitempty"`
	Language				string	`json:"language,omitempty"`
	City					string	`json:"city,omitempty"`
	Province				string	`json:"province,omitempty"`
	AvatarUrl				string	`json:"avatar_url"`
	Country					string	`json:"country,omitempty"`
	SessionKey				string	`form:"session_key" json:"session_key"`
	Gender					uint64	`form:"gender" json:"gender"`
	Uid						uint64	`json:"uid"`
	ThirdSession			string	`json:"third_session,omitempty"`
	Phone					string
	Budget					float64	`json:"budget,omitempty"`
	BgAvatarUrl				string	`json:"bg_avatar_url,omitempty"`
	BonusPoints				uint64	`json:"bonus_points,omitempty"`
	HeaderPosition1			string	`gorm:"column:header_position_1" json:"header_position_1,omitempty"`
	HeaderPosition2			string	`gorm:"column:header_position_2" json:"header_position_2,omitempty"`
	HeaderPosition3			string	`gorm:"column:header_position_3" json:"header_position_3,omitempty"`
	BgAvatarId				uint64	`json:"bg_avatar_id,omitempty"`
	Remind					uint64	`json:"remind,omitempty"`
	HiddenAssetMoney		uint64	`json:"hidden_asset_money"`
	AlreadyLogin			uint64	`json:"already_login"`
}

func (*User) GetFirst() (*User) {
	user := User{}
	db.First(&user)

	return &user
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

func (user User) sessionKey() string {
	key := fmt.Sprintf("@user_%d_session_key@", user.ID)
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