package rest

type EntityAccessor interface {
	ParseEntity(form *RequestForm) (interface{}, error)
	ParseEntityIdentity(form *RequestForm) (interface{}, error)
	ParseEntityCaption(form *RequestForm) (string, error)
	CreateEntity(entity interface{}) error
	ReadEntity(identity interface{}) (interface{}, error)
	UpdateEntity(identity interface{}, entity interface{}) error
	DeleteEntity(identity interface{}) error
	GetEntityByCaption(caption string) (interface{}, error)
}
