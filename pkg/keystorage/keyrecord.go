package keystorage

import "time"

// RawKeyRecord stores information for a Daily Tracing Key. This uses all data
// in string format so we can transfer if quickly without having to encode / decode
type RawKeyRecord struct {
	// ProcessedAt is the timestamp when the observation was processed and stored
	ProcessedAt string
	// DayNumber is the Day Number on which the Daily Tracing Key was used
	DayNumber string
	// DailyTracingKey is the BASE64-encoded Daily Tracing Key
	DailyTracingKey string
}

// KeyRecord stores information for a Daily Tracing Key.
type KeyRecord struct {
	// ProcessedAt is the timestamp when the observation was processed and stored
	ProcessedAt time.Time
	// DayNumber is the Day Number on which the Daily Tracing Key was used
	DayNumber int
	// DailyTracingKey is the BASE64-encoded Daily Tracing Key
	DailyTracingKey string
}
