package error

import "fmt"

type ServiceNameUniqueError struct {
	id string
}

func NewServiceNameUniqueError(id string) ServiceNameUniqueError {
	return ServiceNameUniqueError{id: id}
}

func (s ServiceNameUniqueError) Error() string {
	return fmt.Sprintf("Entity already exists. ID: %s", s.id)
}
