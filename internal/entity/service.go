package entity

type ServiceId string

type Service struct {
	Id                ServiceId
	Title             string
	DurationInMinutes DurationInMinutes
	Description       string
	CostDescription   string
}
