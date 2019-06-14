package service

import (
	"sync"
	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/utils"
)

var Asset = &assetService{mutex: &sync.Mutex{}}

type assetService struct {
	mutex *sync.Mutex
}

// 资产界面的数据
func (srv *assetService)GetWallet() (res map[string]interface{}) {
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
	var list = []map[string]interface{}{}
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