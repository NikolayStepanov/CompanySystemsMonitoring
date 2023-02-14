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

type MMS interface {
	mmsRequest() []domain.MMSData
	checkMMS(value domain.MMSData) bool
	GetResultMMSData() []domain.MMSData
}

type Services struct {
	SMS SMS
	MMS MMS
}

func NewServices(countriesAlphaStorage storages.CountriesAlphaStorager) *Services {
	smsService := NewSMSService(countriesAlphaStorage)
	mmsService := NewMMSService(countriesAlphaStorage)
	return &Services{smsService, mmsService}
}
