package model

import "time"

//CalendarEvent Calendar Event Model
type CalendarEvent struct {
	ID          uint32    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Time        time.Time `json:"time"`
	CalendarID  uint32    `json:"calendarId"`
}
