package queue

import "time"

const (
	EventChoice = "choice"
	EventClick  = "click"
)

//Notification calendar event notification
type Notification struct {
	EventType string    `json:"eventType"`
	BannerID  int64     `json:"bannerId"`
	GroupID   int64     `json:"groupId"`
	Time      time.Time `json:"startTime"`
}
