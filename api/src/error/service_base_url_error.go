package error

type ServiceBaseUrlError struct {
	error string
}

func NewServiceBaseUrlError(error string) ServiceBaseUrlError {
	return ServiceBaseUrlError{error: error}
}

func (s ServiceBaseUrlError) Error() string {
	return s.error
}
