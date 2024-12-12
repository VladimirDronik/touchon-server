package client

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

// Global instance
var I *Client

func New() *Client {
	return &Client{
		client: &fasthttp.Client{
			Name:               "fasthttp",
			ReadTimeout:        5 * time.Second,
			WriteTimeout:       5 * time.Second,
			MaxConnWaitTimeout: 5 * time.Second,
		},
	}
}

type Client struct {
	client           *fasthttp.Client
	ignoreHttpStatus bool
}

func (o *Client) GetIgnoreHttpStatus() bool {
	return o.ignoreHttpStatus
}

func (o *Client) SetIgnoreHttpStatus(v bool) {
	o.ignoreHttpStatus = v
}

func (o *Client) DoRequest(method, uri string, params, headers map[string]string, body interface{}) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Формируем параметры запроса.
	if len(params) > 0 {
		vals := url.Values{}
		for k, v := range params {
			vals.Add(k, v)
		}
		uri += "?" + vals.Encode()
	}

	req.Header.SetMethod(method)
	req.Header.SetRequestURI(uri)

	// Говорим о поддержке сжатия ответа.
	req.Header.Set("Accept-Encoding", "gzip")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	switch v := body.(type) {
	case []byte:
		req.SetBody(v)
	case nil:
	default:
		data, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "http.Client.DoRequest")
		}
		req.SetBody(data)
	}

	if err := o.client.Do(req, resp); err != nil {
		return nil, errors.Wrap(err, "http.Client.DoRequest")
	}

	var respBody []byte
	var err error

	// Если ответ сжат - распаковываем.
	switch string(resp.Header.ContentEncoding()) {
	case "gzip":
		respBody, err = resp.BodyGunzip()
		if err != nil {
			return nil, errors.Wrap(err, "http.Client.DoRequest")
		}

	default:
		respBody = resp.Body()
	}

	if len(respBody) == 0 {
		// Проверяем статус, если тело ответа пустое.
		if resp.StatusCode() != http.StatusOK && !o.ignoreHttpStatus {
			return nil, errors.Wrap(errors.Errorf("HTTP#%d %s", resp.StatusCode(), string(resp.Body())), "http.Client.DoRequest")
		}

		return nil, nil
	}

	// Сначала проверяем ошибку в теле - она более информативна, чем http статус.
	errMsg := &errResp{}
	if err := json.Unmarshal(resp.Body(), errMsg); err == nil && errMsg.Error != "" {
		return nil, errors.Wrap(errors.New(errMsg.Error), "http.Client.DoRequest")
	}

	// Проверяем статус ответа.
	if resp.StatusCode() != http.StatusOK && !o.ignoreHttpStatus {
		return nil, errors.Wrap(errors.Errorf("HTTP#%d %s", resp.StatusCode(), string(resp.Body())), "http.Client.DoRequest")
	}

	r := make([]byte, len(respBody))
	copy(r, respBody)

	return r, nil
}

type errResp struct {
	Error string `json:"error"`
}
