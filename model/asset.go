package model

import (
	"fmt"
	"github.com/yigger/JZ-back/conf"
	"github.com/leekchan/accounting"
)

type Asset struct {
	CommonModel

	Name					string  	`json:"name"`
	Amount					float64		`json:"amount"`
	ParentId				int     	`json:"parent_id"`
	Type					string		`json:"type"`
	Order					int     	`json:"order"`
	IconPath				string     	`json:"icon_path"`
	Lock					int 		`json:"lock"`
	Budget					float64		`json:"budget"`
	Frequent				int			`json:"frequent"`
	Remark					string		`json:"remark"`
	CreatorId				uint			`json:"creator_id"`
}


func (Asset) GetAssetById(id int) *Asset {
	ret := &Asset{}
	if err := db.First(&ret, id).Error; err != nil {
		fmt.Println(err)
		return nil
	}
	return ret
}

func (asset *Asset) IconUrl() string {
	return conf.Host() + asset.IconPath
}

func (asset *Asset) AmountHuman() string {
	ac := accounting.Accounting{Symbol: "", Precision: 2}
	return ac.FormatMoney(asset.Amount)
}
