package model

import (
	"time"
)

type Statement struct {
	CommonModel

	UserId					int  		`json:"user_id"`
	CategoryId				int  		`json:"category_id"`
	AssetId					int  		`json:"asset_id"`
	Amount					float64  	`json:"amount"`
	Type					string      `json:"type"`
	Description				string      `json:"description"`
	Year					int			`json:"year"`
	Month					int
	Day						int
	Time					time.Time
	Residue					float64
	Location				string
	Nation					string
	Province				string
	City					string
	District				string
	Street					string
	TargetAssetId			int
	Title					string
}