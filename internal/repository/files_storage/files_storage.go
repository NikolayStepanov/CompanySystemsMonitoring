package files_storage

import (
	"CompanySystemsMonitoring/internal/repository/files_storage/csv_file"
	"CompanySystemsMonitoring/internal/repository/storages"
)

type FilesStorage struct {
	CountriesAlphaStorage storages.CountriesAlphaStorage
	CountriesCSV          csv_file.CSVFile
}

func NewFileStorage() *FilesStorage {
	return &FilesStorage{}
}

func (C *FilesStorage) LoadingCountries(countries [][]string) error {
	err := error(nil)
	return err
}
