package util

import (
	"github.com/pkg/errors"
	"vpn-web.funcworks.net/model"
)

func HttpVpnGet[T any](url string) (T, error) {
	var zero T
	obj, err := HttpVpnSend[T]("GET", url, nil)
	if err != nil {
		return zero, err
	}
	if obj.Code != 0 {
		return zero, errors.New(obj.Msg)
	}
	return obj.Data, nil
}

func HttpVpnPost[T any](url string, data map[string]any) (T, error) {
	var zero T
	obj, err := HttpVpnSend[T]("POST", url, nil)
	if err != nil {
		return zero, err
	}
	if obj.Code != 0 {
		return zero, errors.New(obj.Msg)
	}
	return obj.Data, nil
}

func HttpVpnSend[T any](method, url string, data map[string]any, optionFuncs ...HttpOptionFunc) (model.Response[T], error) {
	var zero model.Response[T]
	obj, err := HttpSend[model.Response[T]](method, url, data)
	if err != nil {
		return zero, err
	}
	return obj, nil
}
