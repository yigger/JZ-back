package service

import (
	"sync"
)

var Asset = &assetService{mutex: &sync.Mutex{}}

type assetService struct {
	mutex *sync.Mutex
}

func (srv *assetService)GetStatements() (res []map[string]interface{}) {
	
}
