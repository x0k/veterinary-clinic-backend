package appointment

type SampleRateInMinutes int

func (s SampleRateInMinutes) Minutes() int {
	return int(s)
}
