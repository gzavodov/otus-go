package model

type Reference interface {
	GetID() int64
	SetID(int64)

	GetCaption() string
	SetCaption(string)
}

type BaseReference struct {
	ID      int64  `json:"id"`
	Caption string `json:"caption"`
}

func (r *BaseReference) GetID() int64 {
	return r.ID
}

func (r *BaseReference) SetID(id int64) {
	r.ID = id
}

func (r *BaseReference) GetCaption() string {
	return r.Caption
}

func (r *BaseReference) SetCaption(caption string) {
	r.Caption = caption
}
