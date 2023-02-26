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

const CountFieldsVoiceCall = 8

const (
	ConnectionStabilityColumn = 4 + iota
	TTFBColumn
	VoicePurityColumn
	MedianOfCallsTimeColumn
)

type VoiceCallService struct {
	CountriesAlphaStorage storages.CountriesAlphaStorager
}

func NewVoiceCallService(countriesAlphaStorage storages.CountriesAlphaStorager) *VoiceCallService {
	return &VoiceCallService{CountriesAlphaStorage: countriesAlphaStorage}
}

// voiceCallRead read voiceCall data
func (v VoiceCallService) voiceCallRead(path string) []domain.VoiceCallData {
	voiceCallDataResult := []domain.VoiceCallData{}
	file, err := os.Open(path)
	if err != nil {
		log.Println("Cannot open voiceCallData file:", err)
	}
	defer file.Close()
	reader, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Cannot read voiceCallData file:", err)
	}
	lines := strings.Split(string(reader), "\n")
	for _, line := range lines {
		lineSplit := strings.Split(line, ";")
		if len(lineSplit) == CountFieldsVoiceCall {
			if v.checkVoiceCall(lineSplit) {
				connectionStability, _ := strconv.ParseFloat(lineSplit[ConnectionStabilityColumn], 32)
				connectionStability32 := float32(connectionStability)
				TTFB, _ := strconv.Atoi(lineSplit[TTFBColumn])
				voicePurity, _ := strconv.Atoi(lineSplit[VoicePurityColumn])
				medianOfCallsTime, _ := strconv.Atoi(lineSplit[MedianOfCallsTimeColumn])
				voiceCallData := domain.VoiceCallData{
					Country:             v.CountriesAlphaStorage.GetNameCountryFromAlpha(lineSplit[AlphaColumn]),
					Bandwidth:           lineSplit[BandwidthColumn],
					ResponseTime:        lineSplit[ResponseTimeColumn],
					Provider:            lineSplit[ProviderColumn],
					ConnectionStability: connectionStability32,
					TTFB:                TTFB,
					VoicePurity:         voicePurity,
					MedianOfCallsTime:   medianOfCallsTime,
				}
				voiceCallDataResult = append(voiceCallDataResult, voiceCallData)
			}
		}
	}
	return voiceCallDataResult
}

// checkVoiceCall validation VoiceCall data
func (v VoiceCallService) checkVoiceCall(valueLine []string) bool {
	err := error(nil)
	resultValid := true
	if v.CountriesAlphaStorage.ContainsAlpha(valueLine[AlphaColumn]) {
		percentBandwidth := 0
		if percentBandwidth, err = strconv.Atoi(valueLine[BandwidthColumn]); err != nil {
			log.Printf("Value percent bandwidth %v not valid. Error:%s", valueLine, err.Error())
			resultValid = false
		} else if MinBandwidth <= percentBandwidth && percentBandwidth <= MaxBandwidth {
			if _, err = strconv.Atoi(valueLine[ResponseTimeColumn]); err != nil {
				log.Printf("Value responseTime %v not valid. Error:%s", valueLine, err.Error())
			} else {
				if valueLine[ProviderColumn] != VoiceProvidersMap[valueLine[ProviderColumn]] {
					err = fmt.Errorf("not found provider=%s", valueLine[ProviderColumn])
					log.Printf("Value provider %v not valid. Error:%s", valueLine, err.Error())
					resultValid = false
				} else if _, err = strconv.ParseFloat(valueLine[ConnectionStabilityColumn], 32); err != nil {
					log.Printf("Value connection stability %v not valid. Error:%s", valueLine, err.Error())
					resultValid = false
				} else if _, err = strconv.Atoi(valueLine[TTFBColumn]); err != nil {
					log.Printf("Value TTFB %v not valid. Error:%s", valueLine, err.Error())
					resultValid = false
				} else if _, err = strconv.Atoi(valueLine[VoicePurityColumn]); err != nil {
					log.Printf("Value Voice Purity %v not valid. Error:%s", valueLine, err.Error())
					resultValid = false
				} else if _, err = strconv.Atoi(valueLine[MedianOfCallsTimeColumn]); err != nil {
					log.Printf("Value Median Of Calls Time %v not valid. Error:%s", valueLine, err.Error())
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

// GetResultVoiceCallData get result voice call data systems
func (v VoiceCallService) GetResultVoiceCallData(path string) []domain.VoiceCallData {
	resultVoiceCallData := v.voiceCallRead(path)
	sort.Slice(resultVoiceCallData, func(i, j int) bool {
		return resultVoiceCallData[i].Country < resultVoiceCallData[j].Country
	})
	return resultVoiceCallData
}
