package model

import (
	"strconv"
)

type DictInt int

func (d DictInt) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, strconv.Itoa(int(d))), nil
}

func (d *DictInt) UnmarshalJSON(data []byte) error {
	b, val, err := unescapeJsonString(data)
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

type DictBool bool

func (d DictBool) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, strconv.FormatBool(bool(d))), nil
}
func (d *DictBool) UnmarshalJSON(data []byte) error {
	b, val, err := unescapeJsonString(data)
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
