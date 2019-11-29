package endpoint

//Service represents end point service contract
type Service interface {
	GetServiceName() string
	Start() error
	Stop()
}
