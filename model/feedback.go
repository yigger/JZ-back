package model

type Feedback struct {
	CommonModel
	Type		uint64		`json:"type"`
	Content		string		`json:"content"`
	UserId		uint64		`json:"user_id,omitempty"`
}

func (Feedback) Create(feedback *Feedback) {
	db.Create(&feedback)
}