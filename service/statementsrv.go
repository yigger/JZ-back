package service

import (
	"fmt"
	"time"
	"sync"
	"strconv"
	"github.com/mitchellh/mapstructure"
	"github.com/yigger/jodaTime"

	. "github.com/yigger/JZ-back/conf"
	"github.com/yigger/JZ-back/logs"
	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/utils"
)

var Statement = &statementService{mutex: &sync.Mutex{}}

type statementService struct {
	mutex *sync.Mutex
}

func (src *statementService)GetStatements() (res []map[string]interface{}) {
	statements, err := CurrentUser.GetStatements()
	if err != nil {
		logs.Info("err in statement list")
		return nil
	}

	for _, statement := range statements {
		dateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", statement.Date(), time.Local)

		json := map[string]interface{}{
			"id": statement.ID,
			"type": statement.Type,
			"description": statement.Description,
			"title": statement.Title,
			"money": statement.Amount,
			"date": jodaTime.Format("YYYY-MM-dd", dateTime),
			"category": nil,
			"icon_path": nil, // FIXME: icon的路径
			"asset": nil,
			"time": statement.Time(),
			"location": statement.Location,
			"province": statement.Province,
			"city": statement.City,
			"street": statement.Street,
			"month_day": jodaTime.Format("MM-dd", dateTime),
			"timeStr": jodaTime.Format("MM-dd H:m", dateTime),
			"week": utils.WeekMap[dateTime.Weekday().String()],
		}

		var Category model.Category
		category := Category.GetCategoryById(statement.CategoryId)
		if category != nil {
			json["category"] = category.Name
			json["icon_path"] = Conf.Host + category.IconPath
		}
		
		var Asset model.Asset
		asset := Asset.GetAssetById(statement.AssetId)
		if asset != nil {
			json["asset"] = asset.Name
		}

		res = append(res, json)
	}
	
	return
}

func (src *statementService)CreateStatement(params map[string]interface{}) (*model.Statement, error) {
	statementParams := formatStatementParams(params)
	fmt.Println(statementParams)
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
			fmt.Println(err)
		}

		var toAsset model.Asset
		if err := db.Where("creator_id = ? AND id = ?", CurrentUser.ID, params["to"].(float64)).Find(&toAsset).Error; err != nil {
			fmt.Println(err)
		}

		var category model.Category
		if err := db.Where("user_id = ? AND id = ?", CurrentUser.ID, categoryId).Find(&category).Error; err != nil {
			fmt.Println(err)
		}

		var transferCategory model.Category
		if err := db.Where("user_id = ? AND type = 'transfer' AND name = ?", CurrentUser.ID, "转账").Find(&transferCategory).Error; err != nil {
			fmt.Println(err)
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
					 fmt.Println(err)
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
					fmt.Println(err)
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
		var childs []model.Asset
		db.Where("parent_id = ?", asset.ID).Find(&childs)

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
		var childs []model.Category
		db.Where("parent_id = ?", category.ID).Find(&childs)

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
