package app

import (
	"CompanySystemsMonitoring/internal/config"
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

func Run(configPath string) {
	err := error(nil)
	cfg := new(config.Config)
	cfg, err = config.Init(configPath)
	if err != nil {
		log.Fatalln(err)
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	alphaCSV := csv_file.NewCSVFile(cfg.FileStorageApp.RootPath + cfg.FileStorageApp.Alpha)
	countryAlphaStorage := storages.CountriesAlphaStorage{}
	filesStorage := files_storage.NewFileStorage(&countryAlphaStorage, alphaCSV)
	filesStorage.LoadingCountries()
	resultStorage := storages.NewResultDataStorage()
	services := service.NewServices(&countryAlphaStorage)
	result := service.NewResultService(cfg, services, resultStorage)
	handlers := httpDelivery.NewHandler(result)
	//HTTP Server
	srv := server.NewServer(cfg, handlers.Init())
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
