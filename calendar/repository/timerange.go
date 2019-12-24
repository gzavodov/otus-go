package repository

import "time"

//TimeRange represents time range with start and end
type TimeRange struct {
	From time.Time
	To   time.Time
}
