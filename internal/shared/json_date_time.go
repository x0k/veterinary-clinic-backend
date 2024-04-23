package shared

import "time"

// hh:mm
type JsonTime string

// yyyy-mm-dd
type JsonDate string

func NewJsonDate(date string) (JsonDate, error) {
	_, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return "", err
	}
	return JsonDate(date), nil
}

type JsonDateTime string

func JsonDateToGoTime(date JsonDate) (time.Time, error) {
	return time.Parse(time.DateOnly, string(date))
}

func GoTimeToJsonDate(t time.Time) JsonDate {
	return JsonDate(t.Format(time.DateOnly))
}
