package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/domain/common"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type IncidentService struct {
}

func NewIncidentService() *IncidentService {
	return &IncidentService{}
}

// incidentRequest request for incident data
func (i IncidentService) incidentRequest() ([]domain.IncidentData, error) {
	err := error(nil)
	incidentDataResult := []domain.IncidentData{}
	resp := new(http.Response)
	resp, err = http.Get(common.UrlIncidentSystem)
	defer resp.Body.Close()
	if err != nil {
		log.Println("incidentRequest:", err)
		err = fmt.Errorf("incidentRequest error:%w", err)
	} else if resp.StatusCode == http.StatusOK {
		bodyResp := []byte{}
		if bodyResp, err = ioutil.ReadAll(resp.Body); err != nil {
			log.Println(err)
		} else if err = json.Unmarshal(bodyResp, &incidentDataResult); err != nil {
			log.Println(err)
		}
	} else {
		err = fmt.Errorf("incidentRequest error: stausCode = %d", resp.StatusCode)
	}
	return incidentDataResult, err
}

// GetResultIncidentData get result incident data systems
func (i IncidentService) GetResultIncidentData(ctx context.Context) ([]domain.IncidentData, error) {
	err := error(nil)
	resultIncidentData := []domain.IncidentData{}
	select {
	case <-ctx.Done():
		log.Println("cansel: GetResultIncidentData")
		err = fmt.Errorf("incident data not received")
	default:
		resultIncidentData, err = i.incidentRequest()
		sort.Sort(domain.IncidentByStatus{resultIncidentData})
	}
	return resultIncidentData, err
}
