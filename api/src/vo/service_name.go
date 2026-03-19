package vo

import (
	"regexp"
	apierror "transport/api/src/error"
)

const (
	serviceNameMaxLength = 255
	serviceNameMinLength = 8
)

const (
	ServiceNameValidateLevelStrict ValidateLevel = 1
	ServiceNameValidateLevelNone   ValidateLevel = 0
)

type ServiceName struct {
	value                string
	maxLength, minLength int
}

func NewServiceName(serviceName string, validateLevel ValidateLevel) (ServiceName, error) {
	serviceNameObj := ServiceName{
		value:     serviceName,
		maxLength: serviceNameMaxLength,
		minLength: serviceNameMinLength,
	}

	if validateLevel == ServiceNameValidateLevelStrict {
		isValid, err := serviceNameObj.isValidPattern()

		if err != nil {
			return ServiceName{}, err
		}

		if !isValid {
			return ServiceName{}, apierror.NewServiceNameError("invalid pattern service name")
		}

		isValid = serviceNameObj.isValidLength()

		if !isValid {
			return ServiceName{}, apierror.NewServiceNameError("invalid length service name")
		}
	}

	return serviceNameObj, nil
}

func (n ServiceName) isValidPattern() (bool, error) {
	return regexp.MatchString(`^[a-z0-9_-]+$`, n.value)
}

func (n ServiceName) isValidLength() bool {
	return len([]rune(n.value)) <= n.maxLength && len([]rune(n.value)) >= n.minLength
}

func (n ServiceName) Value() string {
	return n.value
}
