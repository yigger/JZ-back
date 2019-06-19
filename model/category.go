package model
import (
	"fmt"
	"github.com/yigger/JZ-back/conf"
	. "github.com/yigger/JZ-back/log"
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

	User					User
}

// CategoryList
type CategoryItem struct {
	ID				uint		`json:"id"`
	Name			string		`json:"name"`
	Order			int			`json:"order"`
	IconPath		string		`json:"icon_path"`
	ParentId		int			`json:"parent_id"`
	Amount			string		`json:"amount"`
	Type			string 		`json:"type"`
}

type SumResult struct {
	Amount float64 //or int ,or some else
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

func (category *Category) GetAmount() float64 {
	var user User
	if err := db.Model(&category).Related(&user).Error; err != nil {
		Log.Errorf("GetAmount err: " + err.Error())
		return 0
	}
	
	var categoryIds []uint
	if category.ParentId == 0 {
		db.Table("categories").Where("parent_id = ?", category.ID).Pluck("ID", &categoryIds)
	} else {
		categoryIds = append(categoryIds, category.ID)
	}

	var sumResult SumResult
	if err := db.Table("statements").Where("user_id = ? and category_id IN (?)", user.ID, categoryIds).Select("sum(amount) as amount").Scan(&sumResult).Error; err != nil {
		Log.Info(err)
		return 0
	}

	return sumResult.Amount
}

func (_ Category) GetParentCategories(user *User, categoryType string) (categories []Category) {
	if err := db.Model(user).Where("parent_id = 0 AND type = ?", categoryType).Association("Categories").Find(&categories).Error; err != nil {
		Log.Errorf(err.Error())
	}
	return
}
