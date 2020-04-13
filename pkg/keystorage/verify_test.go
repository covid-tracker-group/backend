package keystorage

import (
	"testing"
	"time"

	"simplon.biz/corona/pkg/tools"
)

func TestVerifyRecord(t *testing.T) {
	toDayNumber := tools.TimeToDayNumber(time.Now())
	validRecord := KeyRecord{
		DayNumber:       toDayNumber,
		DailyTracingKey: "MDEyMzQ1Njc4OUFCQ0RFRg==",
	}

	if !VerifyRecord(&validRecord) {
		t.Error("Valid record does not verify")
	}

	if VerifyRecord(&KeyRecord{DayNumber: toDayNumber + 1, DailyTracingKey: validRecord.DailyTracingKey}) {
		t.Error("Data from the future is not rejected")
	}

	if VerifyRecord(&KeyRecord{DayNumber: toDayNumber - 15, DailyTracingKey: validRecord.DailyTracingKey}) {
		t.Error("Data from the far past is not rejected")
	}

	if VerifyRecord(&KeyRecord{DayNumber: validRecord.DayNumber, DailyTracingKey: "invalid"}) {
		t.Error("Invalid base64 is not rejected")
	}

	if VerifyRecord(&KeyRecord{DayNumber: validRecord.DayNumber, DailyTracingKey: "c2hvcnQ="}) {
		t.Error("Key with invalid length is not rejected")
	}

}
