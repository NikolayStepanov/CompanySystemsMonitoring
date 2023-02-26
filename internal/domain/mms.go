package domain

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

type MMSDataArray []MMSData

type MMSByCountry struct {
	MMSDataArray
}

type MMSByProvider struct {
	MMSDataArray
}

func (c MMSDataArray) Len() int {
	return len(c)
}

func (c MMSDataArray) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c MMSByCountry) Less(i, j int) bool {
	return c.MMSDataArray[i].Country < c.MMSDataArray[j].Country
}

func (c MMSByProvider) Less(i, j int) bool {
	return c.MMSDataArray[i].Provider < c.MMSDataArray[j].Provider
}
