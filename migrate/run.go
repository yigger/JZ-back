package main

import (
	"github.com/yigger/JZ-back/conf"
	"github.com/yigger/JZ-back/model"
)

func main() {
	conf.LoadConf()
	db := model.ConnectDB()
	// ------  Table Migrate -----------
	// Create User
	db.AutoMigrate(&model.User{})

	// Create Asset
	db.AutoMigrate(&model.Asset{})

	// Create Category
	db.AutoMigrate(&model.Category{})

	// Create UserAsset
	db.AutoMigrate(&model.UserAsset{})

	// Create Statement
	db.AutoMigrate(&model.Statement{})

	// Create Message
	db.AutoMigrate(&model.Message{})

	// Create Feedback
	db.AutoMigrate(&model.Feedback{})
}
