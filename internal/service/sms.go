package service

import (
	"CompanySystemsMonitoring/internal/domain"
	. "CompanySystemsMonitoring/internal/domain/common"
	"CompanySystemsMonitoring/internal/repository/storages"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const CountFieldsSMS = 4

type SMSService struct {
	CountriesAlphaStorage storages.CountriesAlphaStorager
}

func NewSMSService(countriesAlphaStorage storages.CountriesAlphaStorager) *SMSService {
	return &SMSService{CountriesAlphaStorage: countriesAlphaStorage}
}

// smsRead read sms data
func (S SMSService) smsRead(path string) []domain.SMSData {
	smsDataResult := []domain.SMSData{}
	file, err := os.Open(path)
	if err != nil {
		log.Println("Cannot open smsData file:", err)
	}
	defer file.Close()
	reader, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Cannot read smsData file:", err)
	}
	lines := strings.Split(string(reader), "\n")
	for _, line := range lines {
		lineSplit := strings.Split(line, ";")
		if len(lineSplit) == CountFieldsSMS {
			if S.checkSMS(lineSplit) {
				smsData := domain.SMSData{
					Country:      S.CountriesAlphaStorage.GetNameCountryFromAlpha(lineSplit[AlphaColumn]),
					Bandwidth:    lineSplit[BandwidthColumn],
					ResponseTime: lineSplit[ResponseTimeColumn],
					Provider:     lineSplit[ProviderColumn],
				}
				smsDataResult = append(smsDataResult, smsData)
			}
		}
	}
	return smsDataResult
}

// checkSMS validation sms data
func (S SMSService) checkSMS(valueLine []string) bool {
	err := error(nil)
	resultValid := true
	if S.CountriesAlphaStorage.ContainsAlpha(valueLine[AlphaColumn]) {
		percentBandwidth := 0
		if percentBandwidth, err = strconv.Atoi(valueLine[BandwidthColumn]); err != nil {
			log.Printf("Value percent bandwidth %v not valid. Error:%s", valueLine, err.Error())
			resultValid = false
		} else if MinBandwidth <= percentBandwidth && percentBandwidth <= MaxBandwidth {
			if _, err = strconv.Atoi(valueLine[ResponseTimeColumn]); err != nil {
				log.Printf("Value responseTime %v not valid. Error:%s", valueLine, err.Error())
			} else {
				if valueLine[ProviderColumn] != ProvidersMap[valueLine[ProviderColumn]] {
					err = fmt.Errorf("not found provider=%s", valueLine[ProviderColumn])
					log.Printf("Value provider %v not valid. Error:%s", valueLine, err.Error())
					resultValid = false
				}
			}
		} else {
			err = fmt.Errorf("value percent bandwidth should be between 0 and 100")
			log.Printf("Value percent bandwidth %v not valid. Error:%s", valueLine, err.Error())
			resultValid = false
		}
	} else {
		err = fmt.Errorf("no such code alpha-2 in countries storage")
		log.Printf("Value alpha-2 %v not valid. Error:%s", valueLine, err.Error())
		resultValid = false
	}
	return resultValid
}

// GetResultSMSData get result sms systems
func (S SMSService) GetResultSMSData(path string) []domain.SMSData {
	resultSMSData := S.smsRead(path)
	sort.Slice(resultSMSData, func(i, j int) bool {
		return resultSMSData[i].Country < resultSMSData[j].Country
	})
	return resultSMSData
}
