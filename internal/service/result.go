package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"fmt"
	"log"
)

type ResultService struct {
	services *Services
}

type Result interface {
	getResultSetData() (domain.ResultSetT, error)
	GetResultData() domain.ResultT
}

func NewResultService(services *Services) *ResultService {
	return &ResultService{services}
}

func (r ResultService) getResultSetData() (domain.ResultSetT, error) {
	err := error(nil)
	resultSetData := domain.ResultSetT{}
	sms := r.services.SMS.GetResultSMSData("simulator/data/sms.data")
	mms, errMMS := r.services.MMS.GetResultMMSData()
	voiceCallData := r.services.VoiceCall.GetResultVoiceCallData("simulator/data/voice.data")
	email := r.services.Email.GetResultEmailData("simulator/data/email.data")
	billing := r.services.Billing.BillingRead("simulator/data/billing.data")
	support, errSuppotr := r.services.Support.GetResultSupportData()
	incident, errIncident := r.services.Incident.GetResultIncidentData()
	if errMMS != nil || errSuppotr != nil || errIncident != nil {
		log.Println(errMMS, errSuppotr, errIncident)
		err = fmt.Errorf("error on collect data")
	} else {
		resultSetData = domain.ResultSetT{
			SMS:       sms,
			MMS:       mms,
			VoiceCall: voiceCallData,
			Email:     email,
			Billing:   billing,
			Support:   support,
			Incidents: incident,
		}
	}
	return resultSetData, err
}

func (r ResultService) GetResultData() domain.ResultT {
	err := error(nil)
	resultData := domain.ResultT{}
	resultSetData := domain.ResultSetT{}
	if resultSetData, err = r.getResultSetData(); err != nil {
		resultData = domain.ResultT{false, resultSetData, err.Error()}
	} else {
		resultData = domain.ResultT{
			Status: true,
			Data:   resultSetData,
			Error:  "",
		}
	}
	return resultData
}
