package error

type ServiceNameError struct {
	error string
}

func NewServiceNameError(error string) ServiceNameError {
	return ServiceNameError{error: error}
}

func (s ServiceNameError) Error() string {
	return s.error
}
