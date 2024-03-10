package entity

type ServiceId string

type Service struct {
	Id                ServiceId
	Title             string
	DurationInMinutes int
	Description       string
	CostDescription   string
}

type RecordId string

type RecordStatus string

const (
	RecordAwaits RecordStatus = "awaits"
	RecordInWork RecordStatus = "inWork"
)

type Record struct {
	Id             RecordId
	Status         RecordStatus
	DateTimePeriod DateTimePeriod
	UserId         *UserId
}
