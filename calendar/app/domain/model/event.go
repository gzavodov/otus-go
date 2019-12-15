package model

import "time"

//Event Calendar Event Model
type Event struct {
	ID           int64         `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Location     string        `json:"location"`
	StartTime    time.Time     `json:"startTime"`
	EndTime      time.Time     `json:"endTime"`
	NotifyBefore time.Duration `json:"notifyBefore"`
	UserID       int64         `json:"userId"`
	CalendarID   int64         `json:"calendarId"`
}
