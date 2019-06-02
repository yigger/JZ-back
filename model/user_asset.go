package model

type UserAsset struct {
	CommonModel

	Name	string
	Path	string
	Type	string
	ImageableType	string
	ImageableId		int
	Score			int
	Locked			int
	System			int
}