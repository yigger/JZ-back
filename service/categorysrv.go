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

func (CategoryService) GetCategoryRootList(categoryType string) []model.Category {
	db := model.ConnectDB()

	categories := make([]model.Category, 0)
	if err := db.Model(&CurrentUser).
				Where("`categories`.parent_id = 0 AND `categories`.type = ?", categoryType).
				Association("Categories").
				Find(&categories).Error; err != nil {
		Log.Info(err)
	}

	return categories
}

func (CategoryService) GetCategoryHeader(categoryType string) (float64, float64, float64) {
	db := model.ConnectDB()
	monthExpend := &model.SumResult{}
	yearExpend := &model.SumResult{}
	allExpend := &model.SumResult{}
	db.Table("statements").Where("type = ? and user_id = ? and year = ? and month = ?", categoryType, CurrentUser.ID, utils.CurrentYear, utils.CurrentMonth).Select("sum(amount) as amount").Scan(&monthExpend)
	db.Table("statements").Where("type = ? and user_id = ? and year = ?", categoryType, CurrentUser.ID, utils.CurrentYear).Select("sum(amount) as amount").Scan(&yearExpend)
	db.Table("statements").Where("type = ? and user_id = ?", categoryType, CurrentUser.ID).Select("sum(amount) as amount").Scan(&allExpend)
	return monthExpend.Amount, yearExpend.Amount, allExpend.Amount
}