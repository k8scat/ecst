package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request struct {
	baseUrl string
	client  *http.Client
}

func NewRequest(baseUrl string) *Request {
	return &Request{
		baseUrl: baseUrl,
		client:  http.DefaultClient,
	}
}

func (r *Request) Get(path string, params *url.Values, headers map[string]string) (jsonStr string, err error) {
	api := fmt.Sprintf("%s%s", r.baseUrl, path)
	var URL *url.URL
	URL, err = url.Parse(api)
	if err != nil {
		return
	}
	if params != nil {
		URL.RawQuery = params.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, URL.String(), nil)
	if err != nil {
		return
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	var resp *http.Response
	resp, err = r.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	jsonStr = string(b)
	return
}

func (r *Request) Post(path string, payload io.Reader, headers map[string]string) (body string, err error) {
	api := fmt.Sprintf("%s%s", r.baseUrl, path)
	req, err := http.NewRequest(http.MethodPost, api, payload)
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	var resp *http.Response
	resp, err = r.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	body = string(b)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s: %s", resp.Status, body)
	}
	return
}
