package storages

import (
	"sync"
)

const countriesCap = 249

type CountriesAlphaStorage struct {
	sync.Mutex
	countries map[string]string
}

func (C *CountriesAlphaStorage) Len() int {
	return len(C.countries)
}

type CountriesAlphaStorager interface {
	Init()
	Len() int
	InitCountries(countriesMap map[string]string)
	GetNameCountryFromAlpha(codeAlpha string) string
	GetAllCountriesAlpha() map[string]string
	ContainsAlpha(alphaCode string) bool
}

func NewCountryAlphaStorage() *CountriesAlphaStorage {
	return &CountriesAlphaStorage{}
}

func (C *CountriesAlphaStorage) Init() {
	C.countries = make(map[string]string, countriesCap)
}

func (C *CountriesAlphaStorage) InitCountries(countriesMap map[string]string) {
	C.countries = countriesMap
}

func (C *CountriesAlphaStorage) GetNameCountryFromAlpha(codeAlpha string) string {
	return C.countries[codeAlpha]
}

func (C *CountriesAlphaStorage) GetAllCountriesAlpha() map[string]string {
	return C.countries
}

func (C *CountriesAlphaStorage) ContainsAlpha(alphaCode string) bool {
	_, ok := C.countries[alphaCode]
	return ok
}
