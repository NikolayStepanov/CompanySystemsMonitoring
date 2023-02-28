package service

import (
	"CompanySystemsMonitoring/internal/config"
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/repository/storages"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

const storageUpdateTimeSec = 30

type ResultService struct {
	cfg           *config.Config
	Time          time.Time
	services      *Services
	resultStorage storages.ResultStorager
}

type Result interface {
	getResultSetData(ctx context.Context) (domain.ResultSetT, error)
	GetResultData(ctx context.Context) domain.ResultT
}

func NewResultService(cfg *config.Config, services *Services, resultStorage storages.ResultStorager) *ResultService {
	return &ResultService{cfg, time.Now().Add(-storageUpdateTimeSec * time.Second), services, resultStorage}
}

func (r *ResultService) getResultSetData(ctx context.Context) (domain.ResultSetT, error) {
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
		sms = r.resultSMSData(ctx, &wg)
	}()
	go func() {
		defer wg.Done()
		mms, errMMS = r.resultMMSData(ctx, &wg)
	}()
	go func() {
		defer wg.Done()
		voiceCall = r.resultVoiceCallData(ctx, &wg)
	}()
	go func() {
		defer wg.Done()
		email = r.resultEmailData(ctx, &wg)
	}()
	go func() {
		defer wg.Done()
		billing = r.resultBillingData(ctx, &wg)
	}()
	go func() {
		defer wg.Done()
		support, errSupport = r.resultSupportData(ctx, &wg)
	}()
	go func() {
		defer wg.Done()
		incidents, errIncident = r.resultIncidentData(ctx, &wg)
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

func (r *ResultService) resultSMSData(ctx context.Context, wg *sync.WaitGroup) [][]domain.SMSData {
	outResult := make(chan [][]domain.SMSData)
	smsResult := [][]domain.SMSData{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(outResult)
		smsData := r.services.SMS.GetResultSMSData(ctx, r.cfg.FilesStorageAPI.RootPath+r.cfg.FilesStorageAPI.SmsFile)
		outResult <- smsData
	}()
	select {
	case <-ctx.Done():
		log.Println("Cansel: SMS Data not received")
	case smsResult = <-outResult:
		log.Println("SMS Data received successfully")
	}
	return smsResult
}

func (r *ResultService) resultMMSData(ctx context.Context, wg *sync.WaitGroup) ([][]domain.MMSData, error) {
	outResult := make(chan [][]domain.MMSData)
	mmsResult := [][]domain.MMSData{}
	errMMS := error(nil)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(outResult)
		mmsData := [][]domain.MMSData{}
		mmsData, errMMS = r.services.MMS.GetResultMMSData(ctx)
		outResult <- mmsData
	}()
	select {
	case <-ctx.Done():
		log.Println("Cansel: MMS Data not received")
	case mmsResult = <-outResult:
		log.Println("MMS Data received successfully")
	}
	return mmsResult, errMMS
}

func (r *ResultService) resultVoiceCallData(ctx context.Context, wg *sync.WaitGroup) []domain.VoiceCallData {
	outResult := make(chan []domain.VoiceCallData)
	voiceCallResult := []domain.VoiceCallData{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(outResult)
		voiceCallData := r.services.VoiceCall.GetResultVoiceCallData(ctx, r.cfg.FilesStorageAPI.RootPath+r.cfg.FilesStorageAPI.VoiceFile)
		outResult <- voiceCallData
	}()
	select {
	case <-ctx.Done():
		log.Println("Cansel: VoiceCall Data not received")
	case voiceCallResult = <-outResult:
		log.Println("VoiceCall Data received successfully")
	}
	return voiceCallResult
}

func (r *ResultService) resultEmailData(ctx context.Context, wg *sync.WaitGroup) map[string][][]domain.EmailData {
	outResult := make(chan map[string][][]domain.EmailData)
	emailResult := map[string][][]domain.EmailData{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(outResult)
		emailData := r.services.Email.GetResultEmailData(ctx, r.cfg.FilesStorageAPI.RootPath+r.cfg.FilesStorageAPI.EmailFile)
		outResult <- emailData
	}()
	select {
	case <-ctx.Done():
		log.Println("Cansel: Email Data not received")
	case emailResult = <-outResult:
		log.Println("Email Data received successfully")
	}
	return emailResult
}

func (r *ResultService) resultBillingData(ctx context.Context, wg *sync.WaitGroup) domain.BillingData {
	outResult := make(chan domain.BillingData)
	billingResult := domain.BillingData{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(outResult)
		billingData := r.services.Billing.BillingRead(ctx, r.cfg.FilesStorageAPI.RootPath+r.cfg.FilesStorageAPI.BillingFile)
		outResult <- billingData
	}()
	select {
	case <-ctx.Done():
		log.Println("Cansel: Billing Data not received")
	case billingResult = <-outResult:
		log.Println("Billing Data received successfully")
	}
	return billingResult
}

func (r *ResultService) resultSupportData(ctx context.Context, wg *sync.WaitGroup) ([]int, error) {
	outResult := make(chan []int)
	supportResult := []int{}
	errSupport := error(nil)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(outResult)
		supportData := []int{}
		supportData, errSupport = r.services.Support.GetResultSupportData(ctx)
		outResult <- supportData
	}()
	select {
	case <-ctx.Done():
		log.Println("Cansel: Support Data not received")
	case supportResult = <-outResult:
		log.Println("Support Data received successfully")
	}
	return supportResult, errSupport
}

func (r *ResultService) resultIncidentData(ctx context.Context, wg *sync.WaitGroup) ([]domain.IncidentData, error) {
	errIncident := error(nil)
	outResult := make(chan []domain.IncidentData)
	incidentResult := []domain.IncidentData{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(outResult)
		incidentsData := []domain.IncidentData{}
		incidentsData, errIncident = r.services.Incident.GetResultIncidentData(ctx)
		outResult <- incidentsData
	}()
	select {
	case <-ctx.Done():
		log.Println("Cansel: Incident Data not received")
	case incidentResult = <-outResult:
		log.Println("Incident Data received successfully")
	}
	return incidentResult, errIncident
}

func (r *ResultService) GetResultData(ctx context.Context) domain.ResultT {
	err := error(nil)
	resultSetData := domain.ResultSetT{}
	resultData := domain.ResultT{}
	timeNow := time.Now()
	differenceTime := timeNow.Sub(r.Time)
	if differenceTime > time.Second*storageUpdateTimeSec {
		if resultSetData, err = r.getResultSetData(ctx); err != nil {
			resultData = domain.ResultT{
				false,
				resultSetData,
				err.Error(),
			}
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
		resultData = domain.ResultT{
			true,
			resultSetData,
			"",
		}
	}
	return resultData
}
