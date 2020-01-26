package model

//Statistics represents statistics
type Statistics struct {
	BannerID       int64  `json:"bannerId"`
	GroupID        int64  `json:"groupId"`
	NumberOfShows  uint32 `json:"numberOfShows"`
	NumberOfClicks uint32 `json:"NumberOfClicks"`
}
