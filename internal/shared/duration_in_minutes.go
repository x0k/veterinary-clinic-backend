package shared

type DurationInMinutes int

func (d DurationInMinutes) Int() int {
	return int(d)
}
