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
)

const (
	sixtyMinutes      = 60
	speedTickitHour   = 18
	notLoadBorder     = 9
	averageLoadBorder = 16
)

const (
	notLoad = iota + 1
	averageLoad
	overLoad
)

type SupportService struct {
}

// supportRequest request for support data
func (s SupportService) supportRequest() ([]domain.SupportData, error) {
	err := error(nil)
	supportDataResult := []domain.SupportData{}

	resp, err := http.Get(common.UrlSupportSystem)
	defer resp.Body.Close()
	if err != nil {
		log.Println("supportRequest:", err)
		err = fmt.Errorf("supportRequest error:%w", err)
	} else if resp.StatusCode == http.StatusOK {
		bodyResp := []byte{}
		if bodyResp, err = ioutil.ReadAll(resp.Body); err != nil {
			log.Println(err)
		} else if err = json.Unmarshal(bodyResp, &supportDataResult); err != nil {
			log.Println(err)
		}
	} else {
		err = fmt.Errorf("supportRequest error: stausCode = %d", resp.StatusCode)
	}
	return supportDataResult, err
}

// GetResultSupportData get result support data systems
func (s SupportService) GetResultSupportData(ctx context.Context) ([]int, error) {
	err := error(nil)
	result := []int{}
	select {
	case <-ctx.Done():
		log.Printf("cansel: GetResultSupportData")
		err = fmt.Errorf("support data not received")
	default:
		totalTicket := 0
		averageTime := 0
		load := 0
		supportData, err := s.supportRequest()
		if err != nil {
			err = fmt.Errorf("GetResultSupportData error:%w", err)
		} else {
			for _, value := range supportData {
				totalTicket += value.ActiveTickets
			}
			if totalTicket < notLoadBorder {
				load = notLoad
			} else if totalTicket <= averageLoadBorder {
				load = averageLoad
			} else {
				load = overLoad
			}
			timeToRequest := float64(sixtyMinutes) / float64(speedTickitHour)
			averageTime = int(float64(totalTicket) * timeToRequest)
			result = append(result, load)
			result = append(result, averageTime)
		}
	}

	return result, err
}

func NewSupportService() *SupportService {
	return &SupportService{}
}
