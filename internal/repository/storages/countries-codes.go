package storages

import (
	"sync"
)

const countriesCap = 249

type CountriesAlphaStorage struct {
	sync.Mutex
	countries map[string]string
}

func NewCountryAlphaStorage() *CountriesAlphaStorage {
	countriesStorage := &CountriesAlphaStorage{}
	countriesStorage.Init()
	return countriesStorage
}

func (C *CountriesAlphaStorage) Init() {
	C.countries = make(map[string]string, countriesCap)
}

func (C *CountriesAlphaStorage) Len() int {
	return len(C.countries)
}

func (C *CountriesAlphaStorage) InitializationCountries(countriesMap map[string]string) {
	C.countries = countriesMap
}
