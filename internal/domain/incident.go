package domain

const (
	activeStatus = "active"
	closedStatus = "closed"
)

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"`
}

type IncidentDataArray []IncidentData

type IncidentByStatus struct {
	IncidentDataArray
}

func (I IncidentDataArray) Len() int {
	return len(I)
}

func (I IncidentDataArray) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

func (I IncidentByStatus) Less(i, j int) bool {
	return I.IncidentDataArray[i].Status != I.IncidentDataArray[j].Status &&
		I.IncidentDataArray[i].Status == activeStatus
}
