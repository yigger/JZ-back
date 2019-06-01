package model

import (
	// "time"
	"fmt"
	"github.com/leekchan/accounting"
	// "github.com/yigger/JZ-back/logs"
)

type Statement struct {
	CommonModel

	UserId					uint  		`json:"user_id"`
	CategoryId				int  		`json:"category_id"`
	AssetId					int  		`json:"asset_id"`
	Amount					float64  	`json:"amount"`
	Type					string      `json:"type"`
	Description				string      `json:"description"`
	Year					int			`json:"year"`
	Month					int			`json:"month"`
	Day						int			`json:"day"`
	// Time					time.Duration	`gorm:"column:time" json:"time"` // 不支持 Time 类型
	Residue					float64		`json:"residue"`
	Location				string		`json:"location"`
	Nation					string		`json:"nation"`
	Province				string		`json:"province"`
	City					string		`json:"city"`
	District				string		`json:"district"`
	Street					string		`json:"street"`
	TargetAssetId			int			`json:"target_asset_id"`
	Title					string		`json:"title"`
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

func (Statement) Create(statement *Statement) {
	db.Create(statement)
}

func (st *Statement) AmountHuman() string {
	ac := accounting.Accounting{Symbol: "", Precision: 2}
	return ac.FormatMoney(st.Amount)
}
