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

type VoiceCall interface {
	voiceCallRead(path string) []domain.VoiceCallData
	checkVoiceCall(value []string) bool
	GetResultVoiceCallData(path string) []domain.VoiceCallData
}

type Email interface {
	emailRead(path string) []domain.EmailData
	checkEmail(value []string) bool
	GetResultEmailData(path string) []domain.EmailData
}

type Billing interface {
	BillingRead(path string) domain.BillingData
}

type Support interface {
	supportRequest() []domain.SupportData
	GetResultSupportData() []domain.SupportData
}

type Services struct {
	SMS       SMS
	MMS       MMS
	VoiceCall VoiceCall
	Email     Email
	Billing   Billing
	Support   Support
}

func NewServices(countriesAlphaStorage storages.CountriesAlphaStorager) *Services {
	smsService := NewSMSService(countriesAlphaStorage)
	mmsService := NewMMSService(countriesAlphaStorage)
	voiceCall := NewVoiceCallService(countriesAlphaStorage)
	email := NewEmailService(countriesAlphaStorage)
	billing := NewBillingService()
	support := NewSupportService()
	return &Services{smsService, mmsService, voiceCall, email, billing, support}
}
