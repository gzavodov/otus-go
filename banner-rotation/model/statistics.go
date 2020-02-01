package model

//Statistics represents statistics
type Statistics struct {
	BannerID       int64 `json:"bannerId"`
	GroupID        int64 `json:"groupId"`
	NumberOfShows  int64 `json:"numberOfShows"`
	NumberOfClicks int64 `json:"NumberOfClicks"`
}
