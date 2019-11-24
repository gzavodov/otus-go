package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//RequestQuery auxiliary struct for working with request GET query params
type RequestQuery struct {
	Request *http.Request
	query   url.Values
}

func (q *RequestQuery) parse() url.Values {
	if q.query == nil {
		q.query = q.Request.URL.Query()
	}
	return q.query
}

//ParseUint32 parses unit32 parameter from form by specified name
func (q *RequestQuery) ParseUint32(name string, defaultValue uint32) (uint32, error) {
	value := q.parse().Get(name)
	if len(value) == 0 {
		return defaultValue, nil
	}

	result, err := strconv.ParseUint(value, 10, 32)
	return uint32(result), err
}

//ParseDate parses date parameter from query string by specified name
func (q *RequestQuery) ParseDate(name string, defaultValue time.Time) (time.Time, error) {
	value := q.parse().Get(name)
	if len(value) == 0 {
		return defaultValue, nil
	}

	result, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse time form \"%s\" (%w)", value, err)
	}
	return result, nil
}
