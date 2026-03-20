package error

import "fmt"

type NotFoundError struct {
	modelName string
}

func NewNotFoundError(modelName string) NotFoundError {
	return NotFoundError{modelName: modelName}
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("Entity not found: %s", n.modelName)
}
