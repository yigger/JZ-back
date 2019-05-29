package model

import (
	"time"
)

const (
	MESSAGE_SYSTEM = 0
	MESSAGE_NOTICE_BAR = 1
)

type Message struct {
	CommonModel

	FromUserId					int	
	TargetId					int
	TargetType					int
	Content						string
	ContentType					string
	AvatarUrl					string
	AlreadyRead					string
	PageUrl						string
	SubTitle					string
	Date						time.Time
}