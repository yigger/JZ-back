package service

import (
	"github.com/pkg/errors"
	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/utils"
	"strconv"
	"sync"
)

var Asset = &assetService{mutex: &sync.Mutex{}}

type assetService struct {
	mutex *sync.Mutex
}

// 资产界面的数据
func (srv *assetService) GetWallet() (res map[string]interface{}) {
	db := model.ConnectDB()

	// 资产的顶部信息
	header := map[string]interface{}{
		"total_asset": utils.FormatMoney(totalAssetSum()),
		"total_liability": utils.FormatMoney(debtSum()),
		"net_worth": utils.FormatMoney(netWorthSum()),
	}

	var parentAssets []model.Asset
	if err := db.Model(&CurrentUser).Where("parent_id = 0").Association("Assets").Find(&parentAssets).Error; err != nil {
		logger.Info(err)
	}
	// 资产的列表
	var list []map[string]interface{}
	for _, asset := range parentAssets {
		// 子分类列表
		var childsAsset []model.Asset
		db.Where("parent_id = ?", asset.ID).Find(&childsAsset)
		var childs []map[string]interface{}
		for _, child := range childsAsset {
			tmp := map[string]interface{}{
				"id": child.ID,
				"name": child.Name,
				"amount": child.AmountHuman(),
				"icon_path": child.IconUrl(),
			}
			childs = append(childs, tmp)
		}

		// 获取父资产的总和
		var amount float64
		rows, _ := db.Table("assets").Where("creator_id = ? AND parent_id = ?", CurrentUser.ID, asset.ID).Select("sum(amount)").Rows()
		if rows.Next() {
			if err := rows.Scan(&amount); err != nil {
				amount = 0
			}
		}

		var json = map[string]interface{}{
			"name": asset.Name,
			"amount": utils.FormatMoney(amount),
			"childs": childs,
		}
		list = append(list, json)
	}
	res = map[string]interface{}{
		"header": header,
		"list": list,
	}
	return 
}

func (srv *assetService) GetAssetInformation(assetId string) (interface{}, error) {
	asset := &model.Asset{}
	db := model.ConnectDB()
	if err := db.Table("assets").Where("id = ? AND creator_id = ?", assetId, CurrentUser.ID).Find(&asset).Error; asset.ID == 0 || err != nil {
		return nil, err
	}

	incomeResult := &model.SumResult{}
	expendResult := &model.SumResult{}
	stQuery := db.Table("statements").Where("user_id = ? AND asset_id = ?", CurrentUser.ID, asset.ID)
	if err := stQuery.Where("type = 'income'").Select("sum(amount) as amount").Scan(&incomeResult).Error; err != nil {
		return nil, err
	}

	if err := stQuery.Where("type = 'expend'").Select("sum(amount) as amount").Scan(&expendResult).Error; err != nil {
		return nil, err
	}

	//surplus := incomeResult.Amount - expendResult.Amount
	res := map[string]interface{}{
		"name": asset.Name,
		"income": utils.FormatMoney(incomeResult.Amount),
		"expend": utils.FormatMoney(expendResult.Amount),
		"surplus": utils.FormatMoney(asset.Amount),
		"source_surplus": asset.Amount,
	}

	return res, nil
}

func (srv *assetService) GetAssetTimeLine(assetId string) (interface{}, error) {
	asset := &model.Asset{}
	db := model.ConnectDB()
	if err := db.Table("assets").Where("id = ? AND creator_id = ?", assetId, CurrentUser.ID).Find(&asset).Error; asset.ID == 0 || err != nil {
		return nil, err
	}

	var times []model.Statement
	if err := db.Table("statements").Where("user_id = ? AND asset_id = ?", CurrentUser.ID, asset.ID).Group("year, month").Order("year desc, month desc").Find(&times).Error; err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	for _, time := range times {
		//month_statemenets = @statements.where(year: time.year, month: time.month)
		query := db.Table("statements").Where("user_id = ? AND asset_id = ? AND year = ? AND month = ?", CurrentUser.ID, asset.ID, time.Year, time.Month)
		incomeResult := &model.SumResult{}
		expendResult := &model.SumResult{}
		if err := query.Where("type = 'income'").Select("sum(amount) as amount").Scan(&incomeResult).Error; err != nil {
			return nil, err
		}

		if err := query.Where("type = 'expend'").Select("sum(amount) as amount").Scan(&expendResult).Error; err != nil {
			return nil, err
		}

		surplus := incomeResult.Amount - expendResult.Amount
		json := map[string]interface{}{
			"expend_amount": utils.FormatMoney(expendResult.Amount),
			"income_amount": utils.FormatMoney(incomeResult.Amount),
			"surplus": utils.FormatMoney(surplus),
			"year": time.Year,
			"month": time.Month,
			"hidden": 1,
		}
		data = append(data, json)
	}

	return data, nil
}

func (srv *assetService) GetStatementsByAsset(assetId string, year string, month string) ([]map[string]interface{}, error) {
	asset := &model.Asset{}
	db := model.ConnectDB()
	if err := db.Table("assets").Where("id = ? AND creator_id = ?", assetId, CurrentUser.ID).Find(&asset).Error; asset.ID == 0 || err != nil {
		return nil, err
	}

	var statements []model.Statement
	if err := db.Table("statements").
				Where("user_id = ? AND asset_id = ? AND year = ? AND month = ?", CurrentUser.ID, asset.ID, year, month).Order("created_at desc").
				Find(&statements).
				Error; err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	for _, statement := range statements {
		data = append(data, statement.ToHumanJson())
	}
	return data, nil
}

func (srv *assetService) UpdateSurplus(id string, params map[string]interface{}) (bool, error) {
	//fmt.Println(params)
	db := model.ConnectDB()
	assetId, _ := strconv.ParseInt(id, 10, 64)
	var Asset model.Asset
	asset := Asset.GetAssetById(int(assetId))
	if asset == nil || asset.CreatorId != CurrentUser.ID {
		return false, errors.New("无效的用户")
	}
	// update the surplus
	asset.Amount = params["amount"].(float64)
	db.Save(asset)
	return true, nil
}

// 用户总资产
func totalAssetSum() (amount float64) {
	db := model.ConnectDB()
	rows, _ := db.Table("assets").Where("creator_id = ? AND type = 'deposit'", CurrentUser.ID).Select("sum(amount)").Rows()
	if rows.Next() {
		if err := rows.Scan(&amount); err != nil {
			amount = 0
		}
	}
	return 
}

// 总负债
func debtSum() (amount float64) {
	db := model.ConnectDB()
	rows, _ := db.Table("assets").Where("creator_id = ? AND type = 'debt'", CurrentUser.ID).Select("sum(amount)").Rows()
	if rows.Next() {
		if err := rows.Scan(&amount); err != nil {
			amount = 0
		}
	}
	return 
}


// 净资产
func netWorthSum() float64 {
	return totalAssetSum() - debtSum()
}