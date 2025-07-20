package utils

import "time"

func AddHours(hours int) time.Time {
	now := time.Now()
	return now.Add(time.Hour * time.Duration(hours))
}
