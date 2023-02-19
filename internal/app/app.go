package app

import (
	httpDelivery "CompanySystemsMonitoring/internal/delivery/http"
	"CompanySystemsMonitoring/internal/repository/files_storage"
	"CompanySystemsMonitoring/internal/repository/files_storage/csv_file"
	"CompanySystemsMonitoring/internal/repository/storages"
	"CompanySystemsMonitoring/internal/server"
	"CompanySystemsMonitoring/internal/service"
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
)

func Run() {
	err := error(nil)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	log.Println("run")
	alphaCSV := csv_file.NewCSVFile("data_app/alpha-2.csv")
	countryAlphaStorage := storages.CountriesAlphaStorage{}
	filesStorage := files_storage.NewFileStorage(&countryAlphaStorage, alphaCSV)
	filesStorage.LoadingCountries()
	services := service.NewServices(&countryAlphaStorage)
	log.Println("SMS Service:")
	sms := services.SMS.GetResultSMSData("simulator/data/sms.data")
	log.Println(sms[0])
	log.Println(sms[1])
	log.Println("MMS Service:")
	mms := services.MMS.GetResultMMSData()
	log.Println(mms[0])
	log.Println(mms[1])
	log.Println("VoiceCall Service:")
	log.Println(services.VoiceCall.GetResultVoiceCallData("simulator/data/voice.data"))
	log.Println("Email Service:")
	log.Println(services.Email.GetResultEmailData("simulator/data/email.data"))
	log.Println("Billing Service:")
	log.Println(services.Billing.BillingRead("simulator/data/billing.data"))
	log.Println("Support Service:")
	log.Println(services.Support.GetResultSupportData())
	log.Println("Incident Service:")
	log.Println(services.Incident.GetResultIncidentData())
	log.Println("Result Service:")
	result := service.NewResultService(services)
	log.Println(result.GetResultData())
	handlers := httpDelivery.NewHandler()
	//HTTP Server
	srv := server.NewServer(handlers.Init())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = srv.Run(); err != nil {
			log.Printf("error occurred while running http server: %s\n", err.Error())
		}
	}()
	log.Print("Server started")
	<-ctx.Done()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = srv.Stop(context.Background()); err != nil {
			log.Printf("error occured on server shutting down: %s", err.Error())
		}
	}()
	wg.Wait()
}
