package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/domain/common"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	sixtyMinutes    = 60
	speedTickitHour = 18
)

const (
	notLoad = iota + 1
	averageLoad
	overLoad
)

type SupportService struct {
}

// supportRequest request for support data
func (s SupportService) supportRequest() []domain.SupportData {
	err := error(nil)
	supportDataResult := []domain.SupportData{}

	resp, err := http.Get(common.UrlSupportSystem)
	if err != nil {
		log.Println("supportRequest:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyResp := []byte{}
		if bodyResp, err = ioutil.ReadAll(resp.Body); err != nil {
			log.Println(err)
		} else if err = json.Unmarshal(bodyResp, &supportDataResult); err != nil {
			log.Println(err)
		}
	}
	return supportDataResult
}

// GetResultSupportData get result support data systems
func (s SupportService) GetResultSupportData() []int {
	result := make([]int, 2)
	totalTicket := 0
	averageTime := 0
	load := 0
	supportData := s.supportRequest()
	for _, value := range supportData {
		totalTicket += value.ActiveTickets
	}
	if totalTicket < 9 {
		load = notLoad
	} else if totalTicket <= 16 {
		load = averageLoad
	} else {
		load = overLoad
	}
	timeToRequest := float64(sixtyMinutes) / float64(speedTickitHour)
	averageTime = int(float64(totalTicket) * timeToRequest)
	result = append(result, load)
	result = append(result, averageTime)
	return result
}

func NewSupportService() *SupportService {
	return &SupportService{}
}
