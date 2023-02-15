package app

import (
	"CompanySystemsMonitoring/internal/repository/files_storage"
	"CompanySystemsMonitoring/internal/repository/files_storage/csv_file"
	"CompanySystemsMonitoring/internal/repository/storages"
	"CompanySystemsMonitoring/internal/service"
	"fmt"
	"log"
	"os"
	"time"
)

func Run() {
	log.Println("run")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(pwd)
	alphaCSV := csv_file.NewCSVFile("data_app/alpha-2.csv")
	countryAlphaStorage := storages.CountriesAlphaStorage{}
	filesStorage := files_storage.NewFileStorage(&countryAlphaStorage, alphaCSV)
	filesStorage.LoadingCountries()
	services := service.NewServices(&countryAlphaStorage)
	log.Println("SMS Service:")
	log.Println(services.SMS.GetResultSMSData("simulator/data/sms.data"))
	log.Println("MMS Service:")
	log.Println(services.MMS.GetResultMMSData())
	log.Println("VoiceCall Service:")
	log.Println(services.VoiceCall.GetResultVoiceCallData("simulator/data/voice.data"))
	time.Sleep(time.Minute * 20)
}
