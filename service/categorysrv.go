package service

import (
	"errors"
	"sync"

	. "github.com/yigger/JZ-back/log"
	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/utils"
)

var Category = &CategoryService{mutex: &sync.Mutex{}}

type CategoryService struct {
	mutex *sync.Mutex
}

// 获取分类管理的列表
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

// 获取分类管理的顶部数据信息
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

// 根据 CategoryId 获取账单列表
func (CategoryService) GetStatementByCategoryId(categoryId string) ([]map[string]interface{}) {
	db := model.ConnectDB()

	var statements []model.Statement
	db.Table("statements").
		Where("user_id = ? AND category_id = ?", CurrentUser.ID, categoryId).
		Group("year, month").
		Order("statements.year desc, month desc").
		Scan(&statements)

	res := []map[string]interface{}{}
	for _, st := range statements {
		var childStatements []model.Statement
		db.Table("statements").
			Where("user_id = ? AND category_id = ? AND year = ? AND month = ?", CurrentUser.ID, categoryId, st.Year, st.Month).
			Select("statements.*").Find(&childStatements)

		// 结果集
		var childs []interface{}
		for _, statement := range childStatements {
			childs = append(childs, statement.ToHumanJson())
		}
		res = append(res, map[string]interface{}{
			"year": st.Year,
			"month": st.Month,
			"childs": childs,
		})
	}

	return res
}

func (CategoryService) GetParentList(categoryType string) (res []*model.CategoryItem) {
	var Category model.Category
	categories := Category.GetParentCategories(CurrentUser, categoryType)
	if len(categories) != 0 {
		for _, category := range categories {
			res = append(res, &model.CategoryItem{
				ID:		category.ID,
				Name:	category.Name,
				Order:  category.Order,
				IconPath: category.IconPath,
				ParentId: category.ParentId,
				Type: 	  category.Type,
				Amount: utils.FormatMoney(category.GetAmount()),
			})
		}
	}
	return
}

func (CategoryService) GetCategoryById(categoryId string) (item model.CategoryEdit,err error) {
	db := model.ConnectDB()
	var category model.Category
	err = db.Find(&category, categoryId).Error
	if err != nil {
		return
	}

	if category.UserId != CurrentUser.ID {
		err = errors.New("无效的参数")
	} else {
		item = model.CategoryEdit{
			ID: category.ID,
			Name: category.Name,
			Order: category.Order,
			IconPath: category.IconPath,
			ParentId: category.ParentId,
			Type: category.Type,
			ParentName: category.Parent().Name,
		}
	}
	return
}
