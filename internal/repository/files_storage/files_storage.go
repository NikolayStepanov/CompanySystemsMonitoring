package files_storage

import (
	"CompanySystemsMonitoring/internal/repository/files_storage/csv_file"
	"CompanySystemsMonitoring/internal/repository/storages"
	"log"
)

const (
	NameCountryColumn = iota
	AlphaCountryColumn
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
	countriesRows := [][]string{}
	countriesMap := map[string]string{}
	if countriesRows, err = C.CountriesCSV.ReadAll(); err != nil {
		log.Println(err)
	} else {
		countriesMap = make(map[string]string, len(countriesRows))
		for _, countryRow := range countriesRows {
			countriesMap[countryRow[AlphaCountryColumn]] = countryRow[NameCountryColumn]
		}
		C.CountriesAlphaStorage.InitCountries(countriesMap)
	}
	return err
}
