package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/domain/common"
	"encoding/json"
	"fmt"
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
func (i IncidentService) GetResultIncidentData() ([]domain.IncidentData, error) {
	resultIncidentData, err := i.incidentRequest()
	sort.Sort(domain.IncidentByStatus{resultIncidentData})
	return resultIncidentData, err
}
