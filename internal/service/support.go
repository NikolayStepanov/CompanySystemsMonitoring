package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/domain/common"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
func (s SupportService) GetResultSupportData() []domain.SupportData {
	resultSupportData := s.supportRequest()
	return resultSupportData
}

func NewSupportService() *SupportService {
	return &SupportService{}
}
