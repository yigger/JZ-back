package service

import (
	"sync"
	
	"github.com/yigger/JZ-back/logs"
	"github.com/yigger/JZ-back/model"
	"github.com/yigger/jodaTime"
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
			"date": nil,
			"category": nil,
			"icon_path": nil,
			"asset": nil,
			"time": statement.Time(),
			"location": statement.Location,
			"province": statement.Province,
			"city": statement.City,
			"street": statement.Street,
			"month_day": "",
			"timeStr": jodaTime.Format("MM-dd H:m", statement.CreatedAt),
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
