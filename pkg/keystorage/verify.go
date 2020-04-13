package keystorage

import (
	"encoding/base64"
	"time"

	"simplon.biz/corona/pkg/tools"
)

func VerifyRecord(record *KeyRecord) bool {
	toDayNumber := tools.TimeToDayNumber(time.Now())

	// No time travel
	if record.DayNumber > toDayNumber {
		return false
	}

	// Also no old data for which observations will have been expired
	// by apps.
	if record.DayNumber < (toDayNumber - 14) {
		return false
	}

	key, err := base64.StdEncoding.DecodeString(record.DailyTracingKey)
	if err != nil || len(key) != 16 {
		return false
	}

	return true
}
