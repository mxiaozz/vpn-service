package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

var TimeFormat0 = "2006-01-02 15:04:05"
var TimeFormat1 = "2006/1/2 15:04:05"

type DateTime struct {
	time.Time
}

func DateTimeNow() DateTime {
	return DateTime{Time: time.Now()}
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return []byte("null"), nil
	}

	b := make([]byte, 0, len(TimeFormat0)+len(`""`))
	b = append(b, '"')
	b = dt.AppendFormat(b, TimeFormat0)
	b = append(b, '"')
	return b, nil
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	b, val, err := dt.unescapeJsonString(data)
	if !b {
		return err
	}

	var format = TimeFormat0
	if len(val) > 4 && val[4] == '/' {
		format = TimeFormat1
	}

	if t, err := time.ParseInLocation(format, val, time.Local); err != nil {
		return err
	} else {
		*dt = DateTime{Time: t}
		return nil
	}
}

func (dt DateTime) MarshalText() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat0))
	b = dt.AppendFormat(b, TimeFormat0)
	return b, nil
}

func (dt *DateTime) UnmarshalText(data []byte) error {
	var val = string(data)

	var format = TimeFormat0
	if len(val) > 4 && val[4] == '/' {
		format = TimeFormat1
	}

	if t, err := time.Parse(format, val); err != nil {
		return err
	} else {
		*dt = DateTime{Time: t}
		return nil
	}
}

func (dt *DateTime) unescapeJsonString(data []byte) (bool, string, error) {
	if data == nil {
		return false, "", nil
	}

	str := string(data)
	if str == "" || str == "\"\"" || str == "null" || str == "undefined" {
		return false, "", nil
	}

	// TODO(https://go.dev/issue/47353): Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return false, "", errors.New("DateTime.UnmarshalJSON: input is not a JSON string")
	}

	data = data[len(`"`) : len(data)-len(`"`)]
	return true, string(data), nil
}

func (dt DateTime) Value() (driver.Value, error) {
	if dt.IsZero() {
		return nil, nil
	}
	return dt.Time.Format(TimeFormat0), nil
}

func (dt *DateTime) Scan(value any) (err error) {
	if value == nil {
		dt.Time = time.Time{}
		return
	}

	switch v := value.(type) {
	case time.Time:
		dt.Time = v
	case string:
		dt.Time, err = time.ParseInLocation(TimeFormat0, v, time.Local)
	case []byte:
		dt.Time, err = time.ParseInLocation(TimeFormat0, string(v), time.Local)
	case int:
		dt.Time = time.Unix(int64(v), 0).In(time.Local)
	case int64:
		dt.Time = time.Unix(v, 0).In(time.Local)
	default:
		dt.Time = time.Time{}
	}

	if err != nil {
		fmt.Println("DateTime Scan: " + err.Error())
	}

	return nil
}
