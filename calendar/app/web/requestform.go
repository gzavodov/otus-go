package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

//RequestForm auxiliary struct for working with request POST form
type RequestForm struct {
	Request  *http.Request
	isParsed bool
}

func (f *RequestForm) parse() error {
	if f.isParsed {
		return nil
	}

	err := f.Request.ParseForm()
	if err != nil {
		return fmt.Errorf("could not parse request form (%w)", err)
	}

	f.isParsed = true
	return nil
}

//ParseUint32 parses unit32 parameter from form by specified name
func (f *RequestForm) ParseUint32(name string, defaultValue uint32) (uint32, error) {
	err := f.parse()
	if err != nil {
		return 0, err
	}

	value := f.Request.FormValue(name)
	if len(value) == 0 {
		return defaultValue, nil
	}
	result, err := strconv.ParseUint(value, 10, 32)
	return uint32(result), err
}

//ParseEvent parses calendar event from form
func (f *RequestForm) ParseEvent() (*model.Event, error) {
	var err error
	var str string

	err = f.parse()
	if err != nil {
		return nil, err
	}

	title := f.Request.FormValue("Title")
	description := f.Request.FormValue("Description")
	location := f.Request.FormValue("Location")

	startTime := time.Time{}
	str = f.Request.FormValue("StartTime")
	if len(str) > 0 {
		startTime, err = time.Parse(time.RFC3339, str)
		if err != nil {
			return nil, fmt.Errorf("could not parse StartTime form \"%s\" (%w)", str, err)
		}
	}

	endTime := time.Time{}
	str = f.Request.FormValue("EndTime")
	if len(str) > 0 {
		endTime, err = time.Parse(time.RFC3339, str)
		if err != nil {
			return nil, fmt.Errorf("could not parse EndTime form \"%s\" (%w)", str, err)
		}
	}

	userID := uint64(0)
	str = f.Request.FormValue("UserID")
	if len(str) > 0 {
		userID, err = strconv.ParseUint(str, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("could not parse UserID form \"%s\" (%w)", str, err)
		}
	}

	calendarID := uint64(0)
	str = f.Request.FormValue("CalendarID")
	if len(str) > 0 {
		calendarID, err = strconv.ParseUint(str, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("could not parse calendarID form \"%s\" (%w)", str, err)
		}
	}

	return &model.Event{
			Title:       title,
			Description: description,
			Location:    location,
			StartTime:   startTime,
			EndTime:     endTime,
			UserID:      uint32(userID),
			CalendarID:  uint32(calendarID),
		},
		nil
}
