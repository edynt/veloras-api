package utils

import "strconv"

func Int32ToString(id int32) string {
	return strconv.FormatInt(int64(id), 10)
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
