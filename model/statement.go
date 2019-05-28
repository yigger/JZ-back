package model

import (
	// "time"
	// "github.com/yigger/JZ-back/logs"
)

type Statement struct {
	CommonModel

	UserId					int  		`json:"user_id"`
	CategoryId				int  		`json:"category_id"`
	AssetId					int  		`json:"asset_id"`
	Amount					float64  	`json:"amount"`
	Type					string      `json:"type"`
	Description				string      `json:"description"`
	Year					int			`json:"year"`
	Month					int
	Day						int
	// Time					time.Time	`gorm:"-" json:"time"`
	Residue					float64
	Location				string
	Nation					string
	Province				string
	City					string
	District				string
	Street					string
	TargetAssetId			int
	Title					string
}

func (user User) GetStatements() (statements []*Statement, err error) {
	if err = db.Where("user_id = ?", user.ID).Find(&statements).Error; err != nil {
		// logs.Info(err)
	}

	return
}