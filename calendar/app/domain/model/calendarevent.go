package model

import "time"

//CalendarEvent Calendar Event Model
type CalendarEvent struct {
	ID          uint32
	Title       string
	Description string
	Location    string
	Time        time.Time
	CalendarID  uint32
}
