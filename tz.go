package elation

import (
	"fmt"
	"time"
)

var defaultTimezone *time.Location

type timeWithOptionalZone struct {
	time.Time
}

func init() {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		defaultTimezone = time.UTC
	} else {
		defaultTimezone = loc
	}
}

func (t *timeWithOptionalZone) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	if s == "null" {
		t.Time = time.Time{}
		return nil
	}

	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return fmt.Errorf("invalid time string format: %s", s)
	}

	s = s[1 : len(s)-1]

	// If RFC3339 fails, try parsing ISO8601 without zone
	var rfc3339Err error
	t.Time, rfc3339Err = time.Parse(time.RFC3339, s)
	if rfc3339Err == nil {
		return nil
	}

	var noZoneErr error
	parsedTime, errNoZone := time.ParseInLocation("2006-01-02T15:04:05", s, defaultTimezone)
	if errNoZone == nil {
		t.Time = parsedTime
		return nil
	}

	return fmt.Errorf("parsing time: RFC3339 %w or ISO8061 with no zone %w", rfc3339Err, noZoneErr)
}
