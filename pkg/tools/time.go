package tools

import "time"

// TimeToDayTime returns the Day Numnber for a timestamp
func TimeToDayNumber(when time.Time) int {
	return UnixTimeToDayNumber(when.Unix())
}

// UnixTimeToDayNumber returns the Day Numnber for a timestamp
func UnixTimeToDayNumber(when int64) int {
	return int(when / 86400)
}
