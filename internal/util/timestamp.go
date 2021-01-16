package util

import "time"

func GetTimestamp() string {
	dateLayout := "Mon, 02 Jan 2006 15:04:05 MST"
	date := time.Now()
	return date.Format(dateLayout)
}
