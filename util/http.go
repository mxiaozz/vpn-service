package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"vpn-web.funcworks.net/model"
)

type HttpOptionFunc func(*http.Client, *http.Request)

func HttpGet[T any](url string) (*T, error) {
	obj, err := HttpSend[T]("get", url, nil)
	if err != nil {
		return nil, err
	}
	if obj.Code != 0 {
		return nil, errors.New(obj.Msg)
	}
	return &obj.Data, nil
}

func HttpPost[T any](url string, data map[string]any) (*T, error) {
	obj, err := HttpSend[T]("post", url, data)
	if err != nil {
		return nil, err
	}
	if obj.Code != 0 {
		return nil, errors.New(obj.Msg)
	}
	return &obj.Data, nil
}

func HttpSend[T any](method, url string, data map[string]any, optionFuncs ...HttpOptionFunc) (*model.Response[T], error) {
	var req *http.Request
	var err error

	if strings.EqualFold(method, "get") {
		if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
			return nil, err
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
		if req, err = http.NewRequest(http.MethodPost, url, reader); err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 20 * time.Second}

	// 额外配置
	for _, config := range optionFuncs {
		config(client, req)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var obj model.Response[T]
	if err = json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}
