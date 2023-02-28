package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/repository/storages"
	"context"
)

type SMS interface {
	smsRead(path string) []domain.SMSData
	checkSMS(value []string) bool
	GetResultSMSData(ctx context.Context, path string) [][]domain.SMSData
}

type MMS interface {
	mmsRequest() ([]domain.MMSData, error)
	checkMMS(value domain.MMSData) bool
	GetResultMMSData(ctx context.Context) ([][]domain.MMSData, error)
}

type VoiceCall interface {
	voiceCallRead(path string) []domain.VoiceCallData
	checkVoiceCall(value []string) bool
	GetResultVoiceCallData(ctx context.Context, path string) []domain.VoiceCallData
}

type Email interface {
	emailRead(path string) []domain.EmailData
	checkEmail(value []string) bool
	GetResultEmailData(ctx context.Context, path string) map[string][][]domain.EmailData
}

type Billing interface {
	BillingRead(ctx context.Context, path string) domain.BillingData
}

type Support interface {
	supportRequest() ([]domain.SupportData, error)
	GetResultSupportData(ctx context.Context) ([]int, error)
}

type Incident interface {
	incidentRequest() ([]domain.IncidentData, error)
	GetResultIncidentData(ctx context.Context) ([]domain.IncidentData, error)
}
type Services struct {
	SMS       SMS
	MMS       MMS
	VoiceCall VoiceCall
	Email     Email
	Billing   Billing
	Support   Support
	Incident  Incident
}

func NewServices(countriesAlphaStorage storages.CountriesAlphaStorager) *Services {
	smsService := NewSMSService(countriesAlphaStorage)
	mmsService := NewMMSService(countriesAlphaStorage)
	voiceCall := NewVoiceCallService(countriesAlphaStorage)
	email := NewEmailService(countriesAlphaStorage)
	billing := NewBillingService()
	support := NewSupportService()
	incident := NewIncidentService()
	return &Services{smsService, mmsService, voiceCall, email, billing,
		support, incident}
}
