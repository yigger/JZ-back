package controller

import (
	"github.com/labstack/echo"
	"net/http"

	. "github.com/yigger/JZ-back/log"
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

// 获取某分类的所有账单列表
func GetCategoryStatementsAction(c echo.Context) error {
	categoryId := c.FormValue("category_id")
	res := service.Category.GetStatementByCategoryId(categoryId)
	return c.JSON(http.StatusOK, res)
}

// 获取父级分类的所有列表
func GetParentCategoriesAction(c echo.Context) error {
	categoryType := c.FormValue("type")
	res := service.Category.GetParentList(categoryType)
	return c.JSON(http.StatusOK, res)
}

// 创建分类
func CreateCategoryAction(c echo.Context) error {
	res := RenderJson()
	defer c.JSON(http.StatusOK, res)
	params := make(map[string]model.Category)
	if err := c.Bind(&params); err != nil {
		Log.Errorf(err.Error())
		res.Msg = err.Error()
	} else {
		var Category model.Category
		category := params["category"]
		category.UserId = service.CurrentUser.ID
		Category.Create(category)
	}

	return nil
}

// 获取编辑分类时所需的分类详情数据
func GetCategoryDetail(c echo.Context) error {
	categoryId := c.Param("id")
	category, err := service.Category.GetCategoryById(categoryId)
	if err != nil {
		res := RenderJson()
		res.Msg = err.Error()
		return c.JSON(http.StatusOK, res)
	} else {
		return c.JSON(http.StatusOK, category)
	}
}
