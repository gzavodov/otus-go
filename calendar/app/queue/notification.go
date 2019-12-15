package queue

import "time"

//Notification calendar event notification
type Notification struct {
	Title     string    `json:"title"`
	StartTime time.Time `json:"startTime"`
	UserID    int64     `json:"userId"`
}
