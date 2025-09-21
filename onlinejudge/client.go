package onlinejudge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Client interface {
	GetRequest(data ojData) (*http.Request, error)
}
type ojClient struct {
	url    string
	header http.Header
	method string
}
type ojData struct {
	dataType string
	data     any
}

func (client *ojClient) GetRequest(data ojData) (*http.Request, error) {
	var req *http.Request
	var err error

	switch data.dataType {
	case BODY_JSON:
		jsonData, err := json.Marshal(data.data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}

		req, err = http.NewRequest(client.method, client.url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("failed to create JSON request: %w", err)
		}

	case PARAMS:
		params, ok := data.data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("params data must be a map[string]interface{}")
		}

		queryValues := url.Values{}
		for key, value := range params {
			switch v := value.(type) {
			case string:
				queryValues.Add(key, v)
			case int, int64, float64:
				queryValues.Add(key, fmt.Sprintf("%v", v))
			case []string:
				for _, item := range v {
					queryValues.Add(key, item)
				}
			case bool:
				queryValues.Add(key, strconv.FormatBool(v))
			default:
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					queryValues.Add(key, fmt.Sprintf("%v", v))
				} else {
					queryValues.Add(key, string(jsonBytes))
				}
			}
		}

		if client.method == http.MethodGet {
			encodedParams := queryValues.Encode()
			var finalURL string
			if encodedParams != "" {
				finalURL = client.url + "?" + encodedParams
			} else {
				finalURL = client.url
			}
			req, err = http.NewRequest(client.method, finalURL, nil)
		} else {
			req, err = http.NewRequest(client.method, client.url, bytes.NewBufferString(queryValues.Encode()))
			if err == nil {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
		}
		if err != nil {
			return nil, fmt.Errorf("failed to create params request: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported data type: %s", data.dataType)
	}

	for key, values := range client.header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	req.AddCookie(&http.Cookie{
		Name:  "csrftoken",
		Value: req.Header.Get("X-Csrftoken"),
	})
	req.AddCookie(&http.Cookie{
		Name:  "sessionid",
		Value: req.Header.Get("SessionId"),
	})
	return req, nil
}

func NewPushAccountClient(ojOptions OjOptions) Client {
	urlStr := fmt.Sprintf("https://%s/api/admin/user", ojOptions.Origin)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("X-Csrftoken", ojOptions.CsrfToken)
	header.Set("Origin", ojOptions.Origin)
	header.Set("Referer", urlStr)
	header.Set("SessionId", ojOptions.SessionId)
	header.Set("Cookie", fmt.Sprintf("csrftoken=%s;sessionid=%s", ojOptions.CsrfToken, ojOptions.SessionId))

	return &ojClient{
		url:    urlStr,
		header: header,
		method: http.MethodPost,
	}
}
func NewGetRankClient(ojOptions OjOptions) Client {
	urlStr := fmt.Sprintf("https://%s/api/contest_rank", ojOptions.Origin)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("X-Csrftoken", ojOptions.CsrfToken)
	header.Set("Origin", ojOptions.Origin)
	header.Set("Referer", urlStr)
	header.Set("SessionId", ojOptions.SessionId)
	header.Set("Cookie", fmt.Sprintf("csrftoken=%s;sessionid=%s", ojOptions.CsrfToken, ojOptions.SessionId))

	return &ojClient{
		url:    urlStr,
		header: header,
		method: http.MethodGet,
	}
}
