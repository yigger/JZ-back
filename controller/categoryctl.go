package controller

import (
	"net/http"
	"github.com/labstack/echo"

	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/service"
	"github.com/yigger/JZ-back/utils"
)

// 收入/支出分类管理数据
func GetCategoryListAction(c echo.Context) error {
	categoryType := c.FormValue("type")
	parentId := c.FormValue("parent_id")

	categories := service.Category.GetCategoryList(parentId, categoryType)
	var categoryList []*model.CategoryItem
	for _, category := range categories {
		categoryList = append(categoryList, &model.CategoryItem{
			ID:		category.ID,
			Name:	category.Name,
			Order:  category.Order,
			IconPath: category.IconPath,
			ParentId: category.ParentId,
			Type: 	  category.Type,
			Amount: utils.FormatMoney(category.GetAmount()),
		})
	}

	monthSum, yearSum, AllSum := service.Category.GetCategoryHeader(parentId, categoryType)
	res := map[string]interface{}{}
	res["header"] = map[string]interface{}{
		"month": utils.FormatMoney(monthSum),
		"year": utils.FormatMoney(yearSum),
		"all": utils.FormatMoney(AllSum),
	}
	res["categories"] = categoryList

	return c.JSON(http.StatusOK, res)
}

func GetCategoryStatementsAction(c echo.Context) error {
	categoryId := c.FormValue("category_id")
	res := service.Category.GetStatementByCategoryId(categoryId)
	return c.JSON(http.StatusOK, res)
}
