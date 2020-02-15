package test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
)

type Form struct {
}

func (f *Form) Post(url string, data url.Values, handler http.HandlerFunc) (*bytes.Buffer, error) {
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
