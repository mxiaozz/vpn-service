package model

import (
	"errors"
	"strconv"
	"strings"
)

type DictInt int

func (d DictInt) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, strconv.Itoa(int(d))), nil
}

func (d *DictInt) UnmarshalJSON(data []byte) error {
	b, val, err := d.unescapeJsonString(data)
	if !b {
		return err
	}
	if v, err := strconv.Atoi(val); err != nil {
		return err
	} else {
		*d = DictInt(v)
		return nil
	}
}

func (d *DictInt) unescapeJsonString(data []byte) (bool, string, error) {
	if data == nil {
		return false, "", nil
	}

	str := string(data)
	if str == "" || str == "\"\"" || str == "null" || str == "undefined" {
		return false, "", nil
	}

	// TODO(https://go.dev/issue/47353): Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		str = strings.TrimSpace(str)
		for _, c := range str {
			if c < '0' || c > '9' {
				return false, "", errors.New("DictInt.UnmarshalJSON: input is not a JSON string")
			}
		}
		return true, str, nil
	}

	data = data[len(`"`) : len(data)-len(`"`)]
	return true, string(data), nil
}

type DictBool bool

func (d DictBool) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, strconv.FormatBool(bool(d))), nil
}
func (d *DictBool) UnmarshalJSON(data []byte) error {
	b, val, err := d.unescapeJsonString(data)
	if !b {
		return err
	}
	if v, err := strconv.ParseBool(val); err != nil {
		return err
	} else {
		*d = DictBool(v)
		return nil
	}
}

func (d *DictBool) unescapeJsonString(data []byte) (bool, string, error) {
	if data == nil {
		return false, "", nil
	}

	str := string(data)
	if str == "" || str == "\"\"" || str == "null" || str == "undefined" {
		return false, "", nil
	}

	// TODO(https://go.dev/issue/47353): Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		str = strings.TrimSpace(str)
		if strings.EqualFold(str, "true") {
			return true, "true", nil
		}
		if strings.EqualFold(str, "false") {
			return true, "false", nil
		}
		return false, "", errors.New("DictBool.UnmarshalJSON: input is not a JSON string")
	}

	data = data[len(`"`) : len(data)-len(`"`)]
	return true, string(data), nil
}
