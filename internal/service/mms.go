package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/domain/common"
	"CompanySystemsMonitoring/internal/repository/storages"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type MMSService struct {
	CountriesAlphaStorage storages.CountriesAlphaStorager
}

// mmsRequest request for mms data
func (M MMSService) mmsRequest() []domain.MMSData {
	err := error(nil)
	mmsDataResult := []domain.MMSData{}

	resp, err := http.Get(common.UrlMMSSystem)
	if err != nil {
		log.Println("mmsRequest:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyResp := []byte{}
		if bodyResp, err = ioutil.ReadAll(resp.Body); err != nil {
			log.Println(err)
		} else if err = json.Unmarshal(bodyResp, &mmsDataResult); err != nil {
			log.Println(err)
		} else {
			index := 0
			for _, value := range mmsDataResult {
				if M.checkMMS(value) {
					value.Country = M.CountriesAlphaStorage.GetNameCountryFromAlpha(value.Country)
					mmsDataResult[index] = value
					index++
				}
			}
			mmsDataResult = mmsDataResult[:index]
		}
	}
	return mmsDataResult
}

// checkMMS validation mms data
func (M MMSService) checkMMS(value domain.MMSData) bool {
	err := error(nil)
	resultValid := true
	if M.CountriesAlphaStorage.ContainsAlpha(value.Country) {
		percentBandwidth := 0
		if percentBandwidth, err = strconv.Atoi(value.Bandwidth); err != nil {
			log.Printf("Value percent bandwidth %v not valid. Error:%s", value, err.Error())
			resultValid = false
		} else if common.MinBandwidth <= percentBandwidth && percentBandwidth <= common.MaxBandwidth {
			if _, err = strconv.Atoi(value.ResponseTime); err != nil {
				log.Printf("Value responseTime %v not valid. Error:%s", value, err.Error())
			} else {
				if value.Provider != common.ProvidersMap[value.Provider] {
					err = fmt.Errorf("not found provider=%s", value.Provider)
					log.Printf("Value provider %v not valid. Error:%s", value, err.Error())
					resultValid = false
				}
			}
		} else {
			err = fmt.Errorf("value percent bandwidth should be between 0 and 100")
			log.Printf("Value percent bandwidth %v not valid. Error:%s", value, err.Error())
			resultValid = false
		}
	} else {
		err = fmt.Errorf("no such code alpha-2 in countries storage")
		log.Printf("Value alpha-2 %v not valid. Error:%s", value, err.Error())
		resultValid = false
	}
	return resultValid
}

// GetResultMMSData get result mms systems
func (M MMSService) GetResultMMSData() []domain.MMSData {
	resultSMSData := M.mmsRequest()
	sort.Slice(resultSMSData, func(i, j int) bool {
		return resultSMSData[i].Country < resultSMSData[j].Country
	})
	return resultSMSData
}

func NewMMSService(countriesAlphaStorage storages.CountriesAlphaStorager) *MMSService {
	return &MMSService{CountriesAlphaStorage: countriesAlphaStorage}
}
