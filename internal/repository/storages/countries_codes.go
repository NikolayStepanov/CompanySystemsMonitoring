package storages

import (
	"sync"
)

const countriesCap = 249

type CountriesAlphaStorage struct {
	sync.Mutex
	countries      map[string]string
	countriesAlpha map[string]string
}

func (C *CountriesAlphaStorage) Len() int {
	return len(C.countries)
}

type CountriesAlphaStorager interface {
	Init()
	Len() int
	InitCountries(countriesMap map[string]string)
	GetNameCountryFromAlpha(codeAlpha string) string
	GetAlphaFromNameCountry(nameCountry string) string
	GetAllCountriesAlpha() map[string]string
	ContainsAlpha(alphaCode string) bool
}

func NewCountryAlphaStorage() *CountriesAlphaStorage {
	return &CountriesAlphaStorage{}
}

func (C *CountriesAlphaStorage) Init() {
	C.countries = make(map[string]string, countriesCap)
	C.countriesAlpha = make(map[string]string, countriesCap)
}

func (C *CountriesAlphaStorage) InitCountries(countriesMap map[string]string) {
	C.Init()
	C.countries = countriesMap
	for key, value := range countriesMap {
		C.countriesAlpha[value] = key
	}
}

func (C *CountriesAlphaStorage) GetNameCountryFromAlpha(codeAlpha string) string {
	return C.countries[codeAlpha]
}

func (C *CountriesAlphaStorage) GetAlphaFromNameCountry(nameCountry string) string {
	return C.countries[nameCountry]
}

func (C *CountriesAlphaStorage) GetAllCountriesAlpha() map[string]string {
	return C.countries
}

func (C *CountriesAlphaStorage) ContainsAlpha(alphaCode string) bool {
	_, ok := C.countries[alphaCode]
	return ok
}
