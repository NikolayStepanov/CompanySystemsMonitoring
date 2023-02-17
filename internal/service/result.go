package service

import "CompanySystemsMonitoring/internal/domain"

type ResultService struct {
	service *Services
}

type Result interface {
	getResultSetData() domain.ResultSetT
	GetResultData() domain.ResultT
}

func NewResultService(services *Services) *ResultService {
	return &ResultService{services}
}

func (r ResultService) getResultSetData() domain.ResultSetT {
	resultSetData := domain.ResultSetT{}
	return resultSetData
}

func (r ResultService) GetResultData() domain.ResultT {
	resultData := domain.ResultT{}
	return resultData
}
