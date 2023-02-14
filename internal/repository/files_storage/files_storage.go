package files_storage

import (
	"CompanySystemsMonitoring/internal/repository/files_storage/csv_file"
	"CompanySystemsMonitoring/internal/repository/storages"
	"log"
)

type FilesStorage struct {
	CountriesAlphaStorage *storages.CountriesAlphaStorage
	CountriesCSV          csv_file.CSVFiler
}

func NewFileStorage(storage *storages.CountriesAlphaStorage, csvFile csv_file.CSVFiler) *FilesStorage {
	return &FilesStorage{CountriesAlphaStorage: storage, CountriesCSV: csvFile}
}

// LoadingCountries load alpha-2 codes from file into storage
func (C *FilesStorage) LoadingCountries() error {
	err := error(nil)
	var countriesRows [][]string
	var countriesMap map[string]string
	if countriesRows, err = C.CountriesCSV.ReadAll(); err != nil {
		log.Println(err)
	} else {
		countriesMap = make(map[string]string, len(countriesRows))
		for _, countryRow := range countriesRows {
			countriesMap[countryRow[1]] = countryRow[0]
		}
		C.CountriesAlphaStorage.InitCountries(countriesMap)
	}
	return err
}
