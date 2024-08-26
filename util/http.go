package util

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

type HttpOptionFunc func(*http.Client, *http.Request)

func HttpSend[T any](method, url string, data map[string]any, optionFuncs ...HttpOptionFunc) (T, error) {
	var rsp T

	// request
	req, err := newRequest(strings.ToUpper(method), url, data)
	if err != nil {
		return rsp, err
	}
	client := &http.Client{Timeout: 20 * time.Second}

	// 额外配置
	for _, config := range optionFuncs {
		config(client, req)
	}

	// response
	response, err := client.Do(req)
	if err != nil {
		return rsp, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return rsp, err
	}

	// 返回值
	switch any(rsp).(type) {
	case []byte:
		return any(body).(T), nil
	case string:
		return any(string(body)).(T), nil
	case *string:
		s := (*string)(unsafe.Pointer(&body))
		return any(s).(T), nil
	default:
		typ := reflect.TypeOf(rsp)
		if typ.Kind() == reflect.Ptr {
			v := reflect.New(typ.Elem()).Interface()
			rsp = v.(T)
			if err = json.Unmarshal(body, rsp); err != nil {
				return rsp, err
			}
		} else {
			if err = json.Unmarshal(body, &rsp); err != nil {
				return rsp, err
			}
		}
		return rsp, nil
	}
}

func newRequest(method, url string, data map[string]any) (*http.Request, error) {
	if method == http.MethodGet {
		if req, err := http.NewRequest(method, url, nil); err != nil {
			return nil, err
		} else {
			return req, nil
		}
	} else {
		var reader io.Reader
		if data != nil {
			body, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			reader = strings.NewReader(string(body))
		}
		if req, err := http.NewRequest(method, url, reader); err != nil {
			return nil, err
		} else {
			req.Header.Set("Content-Type", "application/json")
			return req, nil
		}
	}
}
