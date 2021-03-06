package model

import (
	"fmt"
	"github.com/yigger/JZ-back/conf"
	. "github.com/yigger/JZ-back/log"
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
	Phone					string	`gorm:"-" json:"phone,omitempty"`
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

	Statements				[]Statement `gorm:"FOREIGNKEY:UserId;ASSOCIATION_FOREIGNKEY:ID"`
	Messages 				[]Message	//`gorm:"foreignkey:TargetId"`
	Assets					[]Asset 	`gorm:"FOREIGNKEY:CreatorId;ASSOCIATION_FOREIGNKEY:ID"`
	Categories				[]Category `gorm:"FOREIGNKEY:UserId;ASSOCIATION_FOREIGNKEY:ID"`
}

func (*User) GetFirst() (*User) {
	user := User{}
	db.First(&user)

	return &user
}

func (*User) GetUserByThirdSession(session string) (*User) {
	user := User{}
	if err := db.Where("third_session = ?", session).First(&user).Error; err != nil {
		Log.Info(err)
		return nil
	}

	return &user
}

// 是否有未读的消息
func (user *User) WaitReadMessage() bool {
	var count int
	if err := db.Table("messages").Where("target_id = ? AND already_read = 0", user.ID).Count(&count).Error; err != nil {
		Log.Info(err)
		return false
	}

	return count > 0
}

// 是否显示未读消息的顶部提示
func (user *User) ShowNoticeBar() bool {
	var count int
	if err := db.Table("messages").Where("target_id = ? AND already_read = 0 AND target_type = ?", user.ID, MESSAGE_NOTICE_BAR).Count(&count).Error; err != nil {
		Log.Info(err)
		return false
	}

	return count > 0
}

// 最新的一条消息
func (user *User) LatestMessage() Message {
	message := Message{}
	if err := db.Where("target_id = ? AND already_read = 0 AND target_type = ?", user.ID, MESSAGE_NOTICE_BAR).Order("created_at desc").Find(&message).Error; err != nil {
		
	}

	return message
}

func (user *User) BgAvatarPath() string {
	userAsset := &UserAsset{}
	if err := db.Find(&userAsset, user.BgAvatarId).Error; err != nil {
		return ""
	}
	path := ""
	// if userAsset.System == 0 {
	path += userAsset.Path
	// } else {
	// 	path += ""
	// }
	return conf.Host() + "/" + path
}

func (user *User) AvatarPath() string {
	return user.AvatarUrl
}

func (user *User) PersistDay() int {
	count := 0
	err := db.Model(&Statement{}).Where("user_id = ?", user.ID).Select("count(distinct year, month, day)").Count(&count).Error
	if err != nil {
		Log.Info(err)
		count = 0
	}
	return count
}

func (user *User) StatementCount() int {
	var count int
	if err := db.Model(&Statement{}).Where("user_id = ?", user.ID).Select("count(*)").Count(&count).Error; err != nil {
		Log.Info(err)
		count = 0
	}
	return count
}

func (User) GetUserByOpenId(openid string) (*User) {
	user := User{}
	if err := db.Where("openid = ?", openid).First(&user).Error; err != nil {
		Log.Info(err)
		return nil
	}

	return &user
}

func (user User) sessionKey() string {
	key := fmt.Sprintf("@user_%d_session_key@", user.ID)
	return key
}

// 缓存中的 session 值
func (user User) CacheSessionVal() (string) {
	cacheVal, err := redisCli.Get(user.sessionKey()).Result()
	if err != nil {
		Log.Info(err)
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
