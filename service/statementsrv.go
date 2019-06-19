package service

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/yigger/JZ-back/model"
	"strconv"
	"sync"
	"time"
)

var Statement = &statementService{mutex: &sync.Mutex{}}

type statementService struct {
	mutex *sync.Mutex
}

func (srv *statementService)GetStatements() (res []map[string]interface{}) {
	db := model.ConnectDB()

	// 筛选七天内的账单数据
	var statements []model.Statement
	beforeDay, _ := time.ParseDuration("-168h") // 7day before
	seventDayBefore := time.Now().Add(beforeDay).Format("2006-01-02 15:04:05")
	if err := db.Model(&CurrentUser).
				 Where("created_at >= ? AND created_at <= ?", seventDayBefore, time.Now()).
				 Order("created_at desc").
				 Association("Statements").
				 Find(&statements).Error; err != nil {
		return nil
	}

	for _, statement := range statements {
		res = append(res, statement.ToHumanJson())
	}

	return
}

func (srv *statementService)CreateStatement(params map[string]interface{}) (*model.Statement, error) {
	statementParams := formatStatementParams(params)
	statement := &model.Statement{}
	err := mapstructure.Decode(statementParams, &statement)
	if err != nil {
		fmt.Println(err)
		return statement, err
	}

	var Statement model.Statement
	statement.UserId = CurrentUser.ID
	Statement.Create(statement)
	return statement, nil
}

func (srv *statementService)UpdateStatement(statement *model.Statement, params map[string]interface{}) (*model.Statement, error) {
	statementParams := formatStatementParams(params)
	err := mapstructure.Decode(statementParams, &statement)
	if err != nil {
		logger.Error(err)
		return statement, err
	}

	logger.Info(statement)
	return statement, nil
}

func formatStatementParams(params map[string]interface{}) (map[string]interface{}) {
	paramsTime := fmt.Sprintf("%s %s", params["date"], params["time"])
	layout := "2006-01-02 15:04:05"
	time, _ := time.Parse(layout, paramsTime)
	amount, _ := strconv.ParseFloat(params["amount"].(string), 64)

	var categoryId int64
	if _, exist := params["category_id"]; exist {
		categoryId, _ = strconv.ParseInt(params["category_id"].(string), 10, 64)
	}
	
	var assetId int64
	if _, exist := params["asset_id"]; exist {
		assetId, _ = strconv.ParseInt(params["asset_id"].(string), 10, 64)
	}
	
	statementParams := map[string]interface{}{
		"CategoryId": categoryId,
      	"AssetId": assetId,
		"Amount": amount,
		"Type": params["type"],
      	"Description": params["description"],
		"Year": time.Year(),
		"Month": time.Month(),
		"Day": time.Day(),
		"CreatedAt": time,
		"Location": params["location"],
        "Nation": params["nation"],
        "Province": params["province"],
        "City": params["city"],
        "District": params["district"],
        "Street": params["street"],
	}

	if params["type"] == "transfer" {
		db := model.ConnectDB()
		var fromAsset model.Asset
		if err := db.Where("creator_id = ? AND id = ?", CurrentUser.ID, params["from"].(float64)).Find(&fromAsset).Error; err != nil {
			logger.Error(err)
		}

		var toAsset model.Asset
		if err := db.Where("creator_id = ? AND id = ?", CurrentUser.ID, params["to"].(float64)).Find(&toAsset).Error; err != nil {
			logger.Error(err)
		}

		var category model.Category
		if err := db.Where("user_id = ? AND id = ?", CurrentUser.ID, categoryId).Find(&category).Error; err != nil {
			logger.Error(err)
		}

		var transferCategory model.Category
		if err := db.Where("user_id = ? AND type = 'transfer' AND name = ?", CurrentUser.ID, "转账").Find(&transferCategory).Error; err != nil {
			logger.Error(err)
		}

		statementParams["CategoryId"] = transferCategory.ID
		statementParams["title"] = fmt.Sprintf("%s -> %s", fromAsset.Name, toAsset.Name)
		statementParams["AssetId"], _ = params["from"].(float64)
		statementParams["TargetAssetId"], _ = params["to"].(float64)
	}

	return statementParams
}

func (*statementService) CategoryFrequentUse(statementType string) ([]model.Category){
	db := model.ConnectDB()

	beforeMin, _ := time.ParseDuration("-30m")
	afterMin, _ := time.ParseDuration("+30m")
	beforeTime := time.Now().Add(beforeMin).Format("15:04:05")
	afterTime := time.Now().Add(afterMin).Format("15:04:05")

	var categories []model.Category
	if err := db.Joins("JOIN statements ON statements.type = ? AND statements.category_id = categories.id", statementType).
				 Where("parent_id > 0 and frequent >= 5 and categories.user_id = ?", CurrentUser.ID).
				 Where(" time(`statements`.created_at) >= ? and time(`statements`.created_at) <= ?", beforeTime, afterTime).
				 Group("categories.id").
				 Order("frequent desc").
				 Limit(3).
				 Find(&categories).Error; err != nil {
					logger.Error(err)
				 }

	return categories
}

func (*statementService) AssetFrequentUse() ([]model.Asset){
	db := model.ConnectDB()
	beforeMin, _ := time.ParseDuration("-30m")
	afterMin, _ := time.ParseDuration("+30m")
	beforeTime := time.Now().Add(beforeMin).Format("15:04:05")
	afterTime := time.Now().Add(afterMin).Format("15:04:05")

	var assets []model.Asset
	if err := db.Joins("JOIN statements ON statements.asset_id = assets.id").
				Where("parent_id > 0 and frequent >= 5 and assets.creator_id = ?", CurrentUser.ID).
				Where(" time(`statements`.created_at) >= ? and time(`statements`.created_at) <= ?", beforeTime, afterTime).
				Group("assets.id").
				Order("frequent desc").
				Limit(3).
				Find(&assets).Error; err != nil {
					logger.Error(err)
				}

	return assets
}

func (*statementService) GetStatementAssets() (res map[string]interface{}) {
	db := model.ConnectDB()

	// 资产列表
	var assets []model.Asset
	if err := db.Model(&CurrentUser).Where("parent_id = 0").Association("Assets").Find(&assets).Error; err != nil {
		panic(err)
	}
	assetRes := []map[string]interface{}{}
	for _, asset := range assets {
		var assetChilds []model.Asset
		db.Where("parent_id = ?", asset.ID).Find(&assetChilds)
		// 组装子类的数据
		var childs []map[string]interface{}
		for _, child := range assetChilds {
			tmp := map[string]interface{}{
				"id": child.ID,
				"name": child.Name,
				"icon_path": child.IconUrl(),
				"amount": child.AmountHuman(),
			}
			childs = append(childs, tmp)
		}

		json := map[string]interface{}{
			"id": asset.ID,
			"name": asset.Name,
			"icon_path": asset.IconPath,
			"childs": childs,
		}
		assetRes = append(assetRes, json)
	}

	// 常用的列表
	var frequents []model.Asset
	if err := db.Model(&CurrentUser).
				 Where("parent_id > 0 and frequent > 5").
				 Order("frequent desc").
				 Limit(10).
				 Association("Assets").
				 Find(&frequents).Error; err != nil {
		panic(err)
	}

	res = map[string]interface{}{
		"categories": assetRes,
		"frequent": frequents,
	}

	return
}


func (*statementService) GetStatementCategories(statementType string) (res map[string]interface{}) {
	db := model.ConnectDB()
	// 分类
	var categories []model.Category
	if err := db.Model(&CurrentUser).Where("parent_id = 0 AND type = ?", statementType).Association("Categories").Find(&categories).Error; err != nil {
		panic(err)
	}
	categoryRes := []map[string]interface{}{}
	for _, category := range categories {
		var childsCategory []model.Category
		db.Where("parent_id = ?", category.ID).Find(&childsCategory)

		// 组装子类的数据
		var childs []map[string]interface{}
		for _, child := range childsCategory {
			tmp := map[string]interface{}{
				"id": child.ID,
				"name": child.Name,
				"icon_path": child.IconUrl(),
			}
			childs = append(childs, tmp)
		}

		json := map[string]interface{}{
			"id": category.ID,
			"name": category.Name,
			"icon_path": category.IconPath,
			"childs": childs,
		}
		categoryRes = append(categoryRes, json)
	}

	// 常用的列表
	var frequents []model.Category
	if err := db.Model(&CurrentUser).
				 Where("parent_id > 0 and frequent > 5 AND type = ?", statementType).
				 Order("frequent desc").
				 Limit(10).
				 Association("Categories").
				 Find(&frequents).Error; err != nil {
		panic(err)
	}

	res = map[string]interface{}{
		"categories": categoryRes,
		"frequent": frequents,
	}

	return
}
