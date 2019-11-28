package repository

//Repository Error Type
const (
	ErrorNone            = 0
	ErrorGeneral         = 1
	ErrorNotFound        = 2
	ErrorInvalidArgument = 3
	ErrorInvalidObject   = 4
)

//NewError creates new repository error
func NewError(code int, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

//Error represents custom repository error
type Error struct {
	code    int
	message string
}

//GetCode retuns repository error code
func (e *Error) GetCode() int {
	return e.code
}

//Error implementation error interface
func (e *Error) Error() string {
	return e.message
}
