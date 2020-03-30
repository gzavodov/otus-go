package httpmonitoring

import (
	"fmt"
)

type RequestInfo struct {
	ServiceName string
	URL         string
	Method      string
	StatusCode  int
}

func (r *RequestInfo) GetStatusCodeLabel() string {
	return fmt.Sprintf("%d", r.StatusCode)
}

func (r *RequestInfo) GetStatusCodeGroupLabel() string {
	return fmt.Sprintf("%d00", r.StatusCode/100)
}
