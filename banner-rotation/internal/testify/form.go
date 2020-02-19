package testify

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
)

type Form struct {
}

func (f *Form) EmulatePost(url string, data url.Values, handler http.HandlerFunc) (*bytes.Buffer, error) {
	encodedData := data.Encode()
	req, err := http.NewRequest(
		"POST",
		url,
		strings.NewReader(encodedData),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))

	rec := httptest.NewRecorder()
	http.HandlerFunc(handler).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		return nil, errors.New(rec.Body.String())
	}

	return rec.Body, nil
}

func (f *Form) Post(url string, data url.Values) ([]byte, error) {
	encodedData := data.Encode()
	req, err := http.NewRequest(
		"POST",
		url,
		strings.NewReader(encodedData),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	return body, nil
}
