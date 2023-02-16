package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/domain/common"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type IncidentService struct {
}

// incidentRequest request for incident data
func (i IncidentService) incidentRequest() []domain.IncidentData {
	err := error(nil)
	incidentDataResult := []domain.IncidentData{}

	resp, err := http.Get(common.UrlIncidentSystem)
	if err != nil {
		log.Println("incidentRequest:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyResp := []byte{}
		if bodyResp, err = ioutil.ReadAll(resp.Body); err != nil {
			log.Println(err)
		} else if err = json.Unmarshal(bodyResp, &incidentDataResult); err != nil {
			log.Println(err)
		}
	}
	return incidentDataResult
}

// GetResultIncidentData get result incident data systems
func (i IncidentService) GetResultIncidentData() []domain.IncidentData {
	resultIncidentData := i.incidentRequest()
	return resultIncidentData
}

func NewIncidentService() *IncidentService {
	return &IncidentService{}
}
