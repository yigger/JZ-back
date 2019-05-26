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

func (src *statementService)GetStatements() (statements []*model.Statement) {
	statements, err := CurrentUser.GetStatements()
	if err != nil {
		logs.Info("err in statement list")
		return 
	}

	return
}
