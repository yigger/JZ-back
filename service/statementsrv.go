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

	categoryId, _ := strconv.ParseInt(params["category_id"].(string), 10, 64)
	assetId, _ := strconv.ParseInt(params["asset_id"].(string), 10, 64)
	amount, _ := strconv.ParseFloat(params["amount"].(string), 64)

	params = map[string]interface{}{
		"CategoryId": categoryId,
      	"AssetId": assetId,
		"Amount": amount,
		"Type": params["type"],
      	"Description": params["description"],
		"Year": time.Year(),
		"Month": time.Month(),
		"Day": time.Day(),
		"CreatedAt": time,
	}

	return params
}
