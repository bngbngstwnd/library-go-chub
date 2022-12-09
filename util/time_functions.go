package util

import "time"

func GetCurrentSecond(input time.Time) int {
	return (input.Hour() * 3600) + input.Minute()*60 + input.Second()
}
