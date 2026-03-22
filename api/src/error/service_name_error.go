package error

type ServiceNameError struct {
	error string
}

func NewServiceNameError(err string) ServiceNameError {
	return ServiceNameError{error: err}
}

func (s ServiceNameError) Error() string {
	return s.error
}
