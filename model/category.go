package model
import (
	"fmt"
	"github.com/yigger/JZ-back/conf"
)
type Category struct {
	CommonModel

	UserId					int  		`json:"user_id"`
	Type					string		`json:"type"`
	Name					string  	`json:"name"`
	ParentId				int     	`json:"parent_id"`
	Order					int     	`json:"order"`
	IconPath				string     	`json:"icon_path"`
	Lock					int 		`json:"lock"`
	Budget					float64		`json:"budget"`
	Frequent				int			`json:"frequent"`
	IsMess					int			`json:"is_mess"`
}

func (Category) GetCategoryById(id int) *Category {
	ret := &Category{}
	if err := db.First(&ret, id).Error; err != nil {
		fmt.Println(err)
		return nil
	}

	return ret
}


func (category *Category) IconUrl() string {
	return conf.Host() + category.IconPath
}
