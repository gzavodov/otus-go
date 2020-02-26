package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gzavodov/otus-go/banner-rotation/model"
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

//ParseInt64 parses int64 parameter from form by specified name
func (f *RequestForm) ParseInt64(name string, defaultValue int64) (int64, error) {
	err := f.parse()
	if err != nil {
		return 0, err
	}

	value := f.Request.FormValue(name)
	if len(value) == 0 {
		return defaultValue, nil
	}
	result, err := strconv.ParseInt(value, 10, 64)
	return result, err
}

//ParseString parses string parameter from form by specified name
func (f *RequestForm) ParseString(name string, defaultValue string) (string, error) {
	err := f.parse()
	if err != nil {
		return "", err
	}

	value := f.Request.FormValue(name)
	if len(value) == 0 {
		return defaultValue, nil
	}
	return value, nil
}

//ParseBanner parses banner from form
func (f *RequestForm) ParseBanner() (*model.Banner, error) {
	if err := f.parse(); err != nil {
		return nil, err
	}

	return &model.Banner{BaseReference: model.BaseReference{Caption: f.Request.FormValue("caption")}}, nil
}

//ParseSlot parses slot from form
func (f *RequestForm) ParseSlot() (*model.Slot, error) {
	if err := f.parse(); err != nil {
		return nil, err
	}

	return &model.Slot{BaseReference: model.BaseReference{Caption: f.Request.FormValue("caption")}}, nil
}

//ParseGroup parses slot from form
func (f *RequestForm) ParseGroup() (*model.Group, error) {
	if err := f.parse(); err != nil {
		return nil, err
	}

	return &model.Group{BaseReference: model.BaseReference{Caption: f.Request.FormValue("caption")}}, nil
}

//ParseBinding parses banner binding from form
func (f *RequestForm) ParseBinding() (*model.Binding, error) {
	if err := f.parse(); err != nil {
		return nil, err
	}

	bannerID, err := f.ParseInt64("bannerId", 0)
	if err != nil {
		return nil, fmt.Errorf("could not parse BannerID: (%w)", err)
	}

	slotID, err := f.ParseInt64("slotId", 0)
	if err != nil {
		return nil, fmt.Errorf("could not parse SlotID: (%w)", err)
	}

	return &model.Binding{BannerID: bannerID, SlotID: slotID}, nil
}
