package model

import (
	"errors"
	"time"
)

var TimeFormat = "2006-01-02 15:04:05"

type DateTime time.Time

func (dt DateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(dt)
	if t.IsZero() {
		return []byte("null"), nil
	}

	b := make([]byte, 0, len(TimeFormat)+len(`""`))
	b = append(b, '"')
	b = t.AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	b, val, err := unescapeJsonString(data)
	if !b {
		return err
	}
	if t, err := time.ParseInLocation(TimeFormat, val, time.Local); err != nil {
		return err
	} else {
		*dt = DateTime(t)
		return nil
	}
}

func (dt DateTime) MarshalText() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat))
	b = time.Time(dt).AppendFormat(b, TimeFormat)
	return b, nil
}

func (dt *DateTime) UnmarshalText(data []byte) error {
	if t, err := time.Parse(TimeFormat, string(data)); err != nil {
		return err
	} else {
		*dt = DateTime(t)
		return nil
	}
}

func (dt DateTime) String() string {
	return time.Time(dt).Format(TimeFormat)
}

func (dt DateTime) Time() time.Time {
	return time.Time(dt)
}

func unescapeJsonString(data []byte) (bool, string, error) {
	if data == nil {
		return false, "", nil
	}

	str := string(data)
	if str == "" || str == "\"\"" || str == "null" || str == "undefined" {
		return false, "", nil
	}

	// TODO(https://go.dev/issue/47353): Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return false, "", errors.New("Time.UnmarshalJSON: input is not a JSON string")
	}

	data = data[len(`"`) : len(data)-len(`"`)]
	return true, string(data), nil
}
