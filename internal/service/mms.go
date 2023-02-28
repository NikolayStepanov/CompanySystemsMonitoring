package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/domain/common"
	"CompanySystemsMonitoring/internal/repository/storages"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type MMSService struct {
	CountriesAlphaStorage storages.CountriesAlphaStorager
}

func NewMMSService(countriesAlphaStorage storages.CountriesAlphaStorager) *MMSService {
	return &MMSService{CountriesAlphaStorage: countriesAlphaStorage}
}

// mmsRequest request for mms data
func (M MMSService) mmsRequest() ([]domain.MMSData, error) {
	err := error(nil)
	mmsDataResult := []domain.MMSData{}

	resp, err := http.Get(common.UrlMMSSystem)
	defer resp.Body.Close()
	if err != nil {
		log.Println("mmsRequest:", err)
		err = fmt.Errorf("mmsRequest error:%w", err)
	} else if resp.StatusCode == http.StatusOK {
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
	} else {
		err = fmt.Errorf("mmsRequest error: stausCode = %d", resp.StatusCode)
	}
	return mmsDataResult, err
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

// GetResultMMSData get result mms data systems
func (M MMSService) GetResultMMSData(ctx context.Context) ([][]domain.MMSData, error) {
	err := error(nil)
	resultMMSData := [][]domain.MMSData{}
	select {
	case <-ctx.Done():
		log.Println("cansel: GetResultMMSData")
		err = fmt.Errorf("mms data not received")
	default:
		mmsData := []domain.MMSData{}
		mmsData, err = M.mmsRequest()
		if err != nil {
			err = fmt.Errorf("GetResultMMSData error:%w", err)
		} else {
			mmsDataSortedByProvider := make([]domain.MMSData, len(mmsData))
			mmsDataSortedByCountry := make([]domain.MMSData, len(mmsData))
			copy(mmsDataSortedByProvider, mmsData)
			copy(mmsDataSortedByCountry, mmsData)
			sort.Sort(domain.MMSByProvider{mmsDataSortedByProvider})
			sort.Sort(domain.MMSByCountry{mmsDataSortedByCountry})
			resultMMSData = append(resultMMSData, mmsDataSortedByProvider)
			resultMMSData = append(resultMMSData, mmsDataSortedByCountry)
		}
	}
	return resultMMSData, err
}
