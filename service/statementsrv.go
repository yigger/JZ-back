package service

import (
	"sync"
	
	"github.com/yigger/JZ-back/logs"
	"github.com/yigger/JZ-back/model"
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
		json := map[string]interface{}{
			"id": statement.ID,
			"type": statement.Type,
			"description": statement.Description,
			"title": statement.Title,
			"money": statement.Amount,
			"date": "暂无",
			"category": "",
			"icon_path": "",
			"asset": "",
			"time": "",
			"location": statement.Location,
			"province": statement.Province,
			"city": statement.City,
			"street": statement.Street,
			"month_day": statement.Year,
			"timeStr": statement.Day,
			"week": "",
		}

		var Category model.Category
		category := Category.GetCategoryById(statement.CategoryId)
		if category != nil {
			json["category"] = category.Name
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
