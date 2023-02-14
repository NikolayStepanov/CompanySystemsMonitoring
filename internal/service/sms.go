package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"CompanySystemsMonitoring/internal/domain/common"
	"CompanySystemsMonitoring/internal/repository/storages"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const countFields = 4
const alphaColumn = 0
const bandwidthColumn = 1
const responseTimeColumn = 2
const providerColumn = 3

type SMSService struct {
	CountriesAlphaStorage storages.CountriesAlphaStorager
}

func NewSMSService(countriesAlphaStorage storages.CountriesAlphaStorager) *SMSService {
	return &SMSService{CountriesAlphaStorage: countriesAlphaStorage}
}

// smsRead read sms data
func (S SMSService) smsRead(path string) []domain.SMSData {
	var smsDataResult []domain.SMSData
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
		if len(lineSplit) == countFields {
			if S.checkSMS(lineSplit) {
				smsData := domain.SMSData{
					Country:      S.CountriesAlphaStorage.GetNameCountryFromAlpha(lineSplit[alphaColumn]),
					Bandwidth:    lineSplit[bandwidthColumn],
					ResponseTime: lineSplit[responseTimeColumn],
					Provider:     lineSplit[providerColumn],
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
	if S.CountriesAlphaStorage.ContainsAlpha(valueLine[alphaColumn]) {
		percentBandwidth := 0
		if percentBandwidth, err = strconv.Atoi(valueLine[bandwidthColumn]); err != nil {
			log.Printf("Value percent bandwidth %v not valid. Error:%s", valueLine, err.Error())
			resultValid = false
		} else if common.MinBandwidth <= percentBandwidth && percentBandwidth <= common.MaxBandwidth {
			if _, err = strconv.Atoi(valueLine[responseTimeColumn]); err != nil {
				log.Printf("Value responseTime %v not valid. Error:%s", valueLine, err.Error())
			} else {
				if valueLine[providerColumn] != common.ProvidersMap[valueLine[providerColumn]] {
					err = fmt.Errorf("provider=%s is absent", valueLine[providerColumn])
					log.Printf("Value provider %v not valid. Error:%s", valueLine, err.Error())
					resultValid = false
				}
			}
		} else {
			err = fmt.Errorf("value percent bandwidth should be between 0 and 100")
			log.Printf("Value percent bandwidth %v not valid. Error:%s", valueLine, err.Error())
			resultValid = false
		}
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
