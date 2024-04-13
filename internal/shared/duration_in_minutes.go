package shared

type DurationInMinutes int

func (d DurationInMinutes) Minutes() int {
	return int(d)
}
