package service

//Service represents service contract
type Service interface {
	GetServiceName() string
	Start() error
	Stop() error
}
