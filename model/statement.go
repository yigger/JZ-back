package model

import (
	// "time"
	"fmt"
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
	// Time					time.Duration	`gorm:"column:time" json:"time"` // 不支持 Time 类型
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
		fmt.Println(err)
	}

	return
}

func (st *Statement) Date() string {
	return fmt.Sprintf("%d-%02d-%02d %s", st.Year, st.Month, st.Day, st.Time())
}

func (st *Statement) Time() string {
	return st.CreatedAt.Format("15:04:05")
}
