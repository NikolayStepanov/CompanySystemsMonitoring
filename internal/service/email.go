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

const CountFieldsEmail = 3

const (
	EmailProviderColumn = 1 + iota
	EmailDeliveryTimeColumn
)

type EmailService struct {
	CountriesAlphaStorage storages.CountriesAlphaStorager
}

func NewEmailService(countriesAlphaStorage storages.CountriesAlphaStorager) *EmailService {
	return &EmailService{CountriesAlphaStorage: countriesAlphaStorage}
}

// emailRead read email data
func (e EmailService) emailRead(path string) []domain.EmailData {
	emailDataResult := []domain.EmailData{}
	file, err := os.Open(path)
	if err != nil {
		log.Println("Cannot open emailData file:", err)
	}
	defer file.Close()
	reader, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Cannot read emailData file:", err)
	}
	lines := strings.Split(string(reader), "\n")
	for _, line := range lines {
		lineSplit := strings.Split(line, ";")
		if len(lineSplit) == CountFieldsEmail {
			if e.checkEmail(lineSplit) {
				deliveryTime, _ := strconv.Atoi(lineSplit[EmailDeliveryTimeColumn])
				emailData := domain.EmailData{
					Country:      e.CountriesAlphaStorage.GetNameCountryFromAlpha(lineSplit[AlphaColumn]),
					Provider:     lineSplit[EmailProviderColumn],
					DeliveryTime: deliveryTime,
				}
				emailDataResult = append(emailDataResult, emailData)
			}
		}
	}
	return emailDataResult
}

// checkEmail validation email data
func (e EmailService) checkEmail(value []string) bool {
	err := error(nil)
	resultValid := true
	if e.CountriesAlphaStorage.ContainsAlpha(value[AlphaColumn]) {
		if value[EmailProviderColumn] != EmailProvidersMap[value[EmailProviderColumn]] {
			err = fmt.Errorf("not found provider=%s", value[EmailProviderColumn])
			log.Printf("Value provider %v not valid. Error:%s", value, err.Error())
			resultValid = false
		} else if _, err = strconv.Atoi(value[EmailDeliveryTimeColumn]); err != nil {
			log.Printf("Value delivery time %v not valid. Error:%s", value, err.Error())
			resultValid = false
		}
	} else {
		err = fmt.Errorf("no such code alpha-2 in countries storage")
		log.Printf("Value alpha-2 %v not valid. Error:%s", value, err.Error())
		resultValid = false
	}
	return resultValid
}

// GetResultEmailData get result email data systems
func (e EmailService) GetResultEmailData(path string) []domain.EmailData {
	resultEmailData := e.emailRead(path)
	sort.Slice(resultEmailData, func(i, j int) bool {
		return resultEmailData[i].Country < resultEmailData[j].Country
	})
	return resultEmailData
}
