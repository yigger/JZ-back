package model

import "fmt"

type Asset struct {
	CommonModel

	UserId					int  		`json:"user_id"`
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
	CreatorId				int			`json:"creator_id"`
}


func (Asset) GetAssetById(id int) *Asset {
	ret := &Asset{}
	if err := db.First(&ret, id).Error; err != nil {
		fmt.Println(err)
		return nil
	}

	return ret
}

