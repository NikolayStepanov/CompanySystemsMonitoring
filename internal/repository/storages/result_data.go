package storages

import (
	"CompanySystemsMonitoring/internal/domain"
	"sync"
)

type ResultDataStorage struct {
	sync.Mutex
	Storage domain.ResultSetT
}

type ResultStorager interface {
	SetResult(result domain.ResultSetT)
	GetResult() domain.ResultSetT
}

func NewResultDataStorage() *ResultDataStorage {
	return &ResultDataStorage{}
}

func (r *ResultDataStorage) SetResult(result domain.ResultSetT) {
	r.Lock()
	r.Storage = result
	defer r.Unlock()
}

func (r *ResultDataStorage) GetResult() domain.ResultSetT {
	return r.Storage
}
