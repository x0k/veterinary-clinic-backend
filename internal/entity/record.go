package entity

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
	Service        Service
}
