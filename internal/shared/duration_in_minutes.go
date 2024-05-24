package shared

type DurationInMinutes int

func NewDurationInMinutes(minutes int) DurationInMinutes {
	return DurationInMinutes(minutes)
}

func (d DurationInMinutes) Int() int {
	return int(d)
}
