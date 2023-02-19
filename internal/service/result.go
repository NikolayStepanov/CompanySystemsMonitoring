package service

import (
	"CompanySystemsMonitoring/internal/domain"
)

type ResultService struct {
	services *Services
}

type Result interface {
	getResultSetData() domain.ResultSetT
	GetResultData() domain.ResultT
}

func NewResultService(services *Services) *ResultService {
	return &ResultService{services}
}

func (r ResultService) getResultSetData() domain.ResultSetT {
	sms := r.services.SMS.GetResultSMSData("simulator/data/sms.data")
	mms := r.services.MMS.GetResultMMSData()
	voiceCallData := r.services.VoiceCall.GetResultVoiceCallData("simulator/data/voice.data")
	email := r.services.Email.GetResultEmailData("simulator/data/email.data")
	billing := r.services.Billing.BillingRead("simulator/data/billing.data")
	support := r.services.Support.GetResultSupportData()
	incident := r.services.Incident.GetResultIncidentData()
	resultSetData := domain.ResultSetT{
		SMS:       sms,
		MMS:       mms,
		VoiceCall: voiceCallData,
		Email:     email,
		Billing:   billing,
		Support:   support,
		Incidents: incident,
	}
	return resultSetData
}

func (r ResultService) GetResultData() domain.ResultT {
	resultData := domain.ResultT{
		Data: r.getResultSetData(),
	}
	return resultData
}
