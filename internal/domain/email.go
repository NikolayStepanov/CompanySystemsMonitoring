package domain

type EmailData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	DeliveryTime int    `json:"delivery_time"`
}

type EmailDataArray []EmailData

type EmailByDeliveryDescending struct {
	EmailDataArray
}

type EmailByDeliveryAscending struct {
	EmailDataArray
}

func (e EmailDataArray) Len() int {
	return len(e)
}

func (e EmailDataArray) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e EmailByDeliveryDescending) Less(i, j int) bool {
	return e.EmailDataArray[i].DeliveryTime < e.EmailDataArray[j].DeliveryTime
}

func (e EmailByDeliveryAscending) Less(i, j int) bool {
	return e.EmailDataArray[i].DeliveryTime > e.EmailDataArray[j].DeliveryTime
}
