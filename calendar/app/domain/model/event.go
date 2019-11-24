package model

import "time"

//Event Calendar Event Model
type Event struct {
	ID          uint32    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	UserID      uint32    `json:"userId"`
	CalendarID  uint32    `json:"calendarId"`
}
