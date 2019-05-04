package models

import "time"

// FlagTime time format used in flag parse.
type FlagTime time.Time

// String get time string.
func (t *FlagTime) String() string {
	return time.Time(*t).Format("2006-01-02 15:04:05")
}

// Set set by time string.
func (t *FlagTime) Set(value string) error {
	parsed, err := time.Parse("2006-01-02 15:04:05", value)

	if err == nil {
		*t = FlagTime(parsed)
	}

	return err
}

// Type get format type string.
func (t *FlagTime) Type() string {
	return "FlagTime"
}

// GetTime get Time format
func (t *FlagTime) GetTime() time.Time {
	return time.Time(*t)
}

func genEmptyTime() FlagTime {
	var emptyTime, _ = time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")

	return FlagTime(emptyTime)
}

// EmptyTime 0001-01-01 00:00:00
var EmptyTime = genEmptyTime()

// JavaTime java timestamp format
type JavaTime int64

// MarshalCSV marshal java time to csv string.
func (t JavaTime) MarshalCSV() string {
	time := time.Unix(0, int64(t)*int64(time.Millisecond))
	return time.Format("2006-01-02 15:04:05.000 Z0700 MST")
}

// UnmarshalCSV unmarshal csv string to java time
func (t JavaTime) UnmarshalCSV(value string) error {
	time, err := time.Parse("2006-01-02 15:04:05.000 Z0700 MST", value)

	if err != nil {
		return err
	}

	t = JavaTime(time.UnixNano() / 1000)
	return nil
}
