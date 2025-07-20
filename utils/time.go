package utils

import "time"

func AddHours(hours int) int64 {
	now := time.Now()
	return now.Add(time.Hour * time.Duration(hours)).Unix()
}

func GetNowUnix() int64 {
	return time.Now().Unix()
}
