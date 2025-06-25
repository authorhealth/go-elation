package elation

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

var defaultTimezone *time.Location

func init() {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		slog.Error("could not load location America/Los_Angeles, falling back to UTC as default timezone", "error", err)
		defaultTimezone = time.UTC
	} else {
		defaultTimezone = loc
	}
}

type TimeWithOptionalZone time.Time

func (t TimeWithOptionalZone) Time() time.Time {
	return time.Time(t)
}

func (t *TimeWithOptionalZone) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}

func (t *TimeWithOptionalZone) UnmarshalJSON(b []byte) error {
	var parsedTime time.Time
	rfc3339Err := parsedTime.UnmarshalJSON(b)
	if rfc3339Err != nil {
		// If RFC3339 parsing fails, try parsing as ISO8601 without zone
		var rawTime string
		err := json.Unmarshal(b, &rawTime)
		if err != nil {
			return fmt.Errorf("parsing raw time as JSON string: %w", err)
		}

		var noZoneErr error
		parsedTime, noZoneErr = time.ParseInLocation("2006-01-02T15:04:05", rawTime, defaultTimezone)
		if noZoneErr != nil {
			return fmt.Errorf("parsing time: RFC3339 %w or ISO8061 with no zone %w", rfc3339Err, noZoneErr)
		}
	}

	*t = TimeWithOptionalZone(parsedTime)
	return nil
}
