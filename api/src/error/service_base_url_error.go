package error

type ServiceBaseURLError struct {
	error string
}

func NewServiceBaseURLError(err string) ServiceBaseURLError {
	return ServiceBaseURLError{error: err}
}

func (s ServiceBaseURLError) Error() string {
	return s.error
}
