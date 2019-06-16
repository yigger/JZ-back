package controller

import (
	"net/http"
	"github.com/labstack/echo"

	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/service"
	"github.com/yigger/JZ-back/utils"
)

func GetCategoryList(c echo.Context) error {
	categoryType := c.FormValue("type")
	categories := service.Category.GetCategoryRootList(categoryType)
	var categoryList []*model.CategoryItem
	for _, category := range categories {
		categoryList = append(categoryList, &model.CategoryItem{
			ID:		category.ID,
			Name:	category.Name,
			Order:  category.Order,
			IconPath: category.IconPath,
			ParentId: category.ParentId,
			Amount: utils.FormatMoney(category.GetAmount()),
		})
	}

	monthSum, yearSum, AllSum := service.Category.GetCategoryHeader(categoryType)
	res := map[string]interface{}{}
	res["header"] = map[string]interface{}{
		"month": utils.FormatMoney(monthSum),
		"year": utils.FormatMoney(yearSum),
		"all": utils.FormatMoney(AllSum),
	}
	res["categories"] = categoryList

	return c.JSON(http.StatusOK, res)
}