package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/repository/storages"
	"fmt"
	"log"
	"sync"
	"time"
)

const storageUpdateTimeSec = 30

type ResultService struct {
	Time          time.Time
	services      *Services
	resultStorage storages.ResultStorager
}

type Result interface {
	getResultSetData() (domain.ResultSetT, error)
	GetResultData() domain.ResultT
}

func NewResultService(services *Services, resultStorage storages.ResultStorager) *ResultService {
	return &ResultService{time.Now().Add(-storageUpdateTimeSec * time.Second), services, resultStorage}
}

func (r *ResultService) getResultSetData() (domain.ResultSetT, error) {
	err := error(nil)
	resultSetData := domain.ResultSetT{}
	var (
		wg          sync.WaitGroup
		errMMS      error
		errSupport  error
		errIncident error
		sms         [][]domain.SMSData
		mms         [][]domain.MMSData
		voiceCall   []domain.VoiceCallData
		email       map[string][][]domain.EmailData
		billing     domain.BillingData
		support     []int
		incidents   []domain.IncidentData
	)
	wg.Add(7)
	go func() {
		defer wg.Done()
		sms = r.services.SMS.GetResultSMSData("simulator/data/sms.data")
	}()
	go func() {
		defer wg.Done()
		mms, errMMS = r.services.MMS.GetResultMMSData()
	}()
	go func() {
		defer wg.Done()
		voiceCall = r.services.VoiceCall.GetResultVoiceCallData("simulator/data/voice.data")
	}()
	go func() {
		defer wg.Done()
		email = r.services.Email.GetResultEmailData("simulator/data/email.data")
	}()
	go func() {
		defer wg.Done()
		billing = r.services.Billing.BillingRead("simulator/data/billing.data")
	}()
	go func() {
		defer wg.Done()
		support, errSupport = r.services.Support.GetResultSupportData()
	}()
	go func() {
		defer wg.Done()
		incidents, errIncident = r.services.Incident.GetResultIncidentData()
	}()
	wg.Wait()
	if errMMS != nil || errSupport != nil || errIncident != nil {
		log.Println(errMMS, errSupport, errIncident)
		err = fmt.Errorf("error on collect data")
	} else {
		resultSetData = domain.ResultSetT{
			SMS:       sms,
			MMS:       mms,
			VoiceCall: voiceCall,
			Email:     email,
			Billing:   billing,
			Support:   support,
			Incidents: incidents,
		}
	}
	return resultSetData, err
}

func (r *ResultService) GetResultData() domain.ResultT {
	err := error(nil)
	resultData := domain.ResultT{}
	resultSetData := domain.ResultSetT{}
	timeNow := time.Now()
	differenceTime := timeNow.Sub(r.Time)
	log.Println(differenceTime)
	if differenceTime > time.Second*storageUpdateTimeSec {
		if resultSetData, err = r.getResultSetData(); err != nil {
			resultData = domain.ResultT{false, resultSetData, err.Error()}
		} else {
			resultData = domain.ResultT{
				Status: true,
				Data:   resultSetData,
				Error:  "",
			}
		}
		r.resultStorage.SetResult(resultSetData)
		r.Time = time.Now()
	} else {
		resultSetData = r.resultStorage.GetResult()
		resultData = domain.ResultT{true, resultSetData, ""}
	}
	return resultData
}
