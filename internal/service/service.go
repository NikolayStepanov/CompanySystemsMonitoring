package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/repository/storages"
)

type SMS interface {
	smsRead(path string) []domain.SMSData
	checkSMS(value []string) bool
	GetResultSMSData(path string) []domain.SMSData
}

type Services struct {
	SMS SMS
}

func NewServices(countriesAlphaStorage storages.CountriesAlphaStorager) *Services {
	smsService := NewSMSService(countriesAlphaStorage)

	return &Services{smsService}
}
