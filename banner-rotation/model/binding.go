package model

//Binding represents relationship between Banner and Slot
type Binding struct {
	ID       int64 `json:"id"`
	BannerID int64 `json:"bannerId"`
	SlotID   int64 `json:"groupId"`
}
