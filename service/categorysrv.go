package service

import (
	"sync"

	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/utils"
	. "github.com/yigger/JZ-back/log"
)

var Category = &CategoryService{mutex: &sync.Mutex{}}

type CategoryService struct {
	mutex *sync.Mutex
}

func (CategoryService) GetCategoryList(parentId string, categoryType string) []model.Category {
	db := model.ConnectDB()
	categories := db.Model(&CurrentUser).Where("`categories`.parent_id = ?", parentId)

	if categoryType != "" {
		categories = categories.Where("`categories`.type = ?", categoryType)
	}

	res := make([]model.Category, 0)
	if err := categories.Association("Categories").Find(&res).Error; err != nil {
		Log.Info(err)
	}

	return res
}

func (CategoryService) GetCategoryHeader(parentId string, categoryType string) (float64, float64, float64) {
	db := model.ConnectDB()
	monthExpend := &model.SumResult{}
	yearExpend := &model.SumResult{}
	allExpend := &model.SumResult{}

	statements := db.Table("statements").Where("user_id = ?", CurrentUser.ID)
	if categoryType != "" {
		statements = statements.Where("type = ?", categoryType)
	} else {
		categoryIds := []int{}
		db.Model(&model.Category{}).Where("user_id = ? AND parent_id = ?", CurrentUser.ID, parentId).Pluck("ID", &categoryIds)
		statements = statements.Where("category_id IN (?)", categoryIds)
	}

	statements.Where("year = ? and month = ?", utils.CurrentYear, utils.CurrentMonth).Select("sum(amount) as amount").Scan(&monthExpend)
	statements.Where("year = ?", utils.CurrentYear).Select("sum(amount) as amount").Scan(&yearExpend)
	statements.Select("sum(amount) as amount").Scan(&allExpend)
	return monthExpend.Amount, yearExpend.Amount, allExpend.Amount
}
