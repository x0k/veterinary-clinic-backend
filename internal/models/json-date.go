package models

import "time"

// hh:mm
type JsonTime string

// yyyy-mm-dd
type JsonDate string

type JsonDateTime string

func JsonDateToGoTime(date JsonDate) (time.Time, error) {
	return time.Parse(time.DateOnly, string(date))
}

func JsonDateTimeToGoTime(date JsonDateTime) (time.Time, error) {
	return time.Parse(time.DateTime, string(date))
}

func GoTimeToJsonDate(t time.Time) JsonDate {
	return JsonDate(t.Format(time.DateOnly))
}
