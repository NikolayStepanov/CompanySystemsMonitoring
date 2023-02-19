package domain

type SMSData struct {
	Country      string `json:"country"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

type SMSDataArray []SMSData

type SMSByCountry struct {
	SMSDataArray
}

type SMSByProvider struct {
	SMSDataArray
}

func (c SMSDataArray) Len() int {
	return len(c)
}

func (c SMSDataArray) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c SMSByCountry) Less(i, j int) bool {
	return c.SMSDataArray[i].Country < c.SMSDataArray[j].Country
}

func (c SMSByProvider) Less(i, j int) bool {
	return c.SMSDataArray[i].Provider < c.SMSDataArray[j].Provider
}
