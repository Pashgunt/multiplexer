package vo

import (
	"regexp"
	apierror "transport/api/src/error"
)

const (
	serviceBaseUrlMaxLength = 255
	serviceBaseUrlMinLength = 8
)

type BaseUrl struct {
	value                string
	maxLength, minLength int
}

func NewBaseUrl(baseUrl string) (BaseUrl, error) {
	baseUrlObj := BaseUrl{
		value:     baseUrl,
		maxLength: serviceBaseUrlMaxLength,
		minLength: serviceBaseUrlMinLength,
	}

	isValid, err := baseUrlObj.isValidPattern()

	if err != nil {
		return BaseUrl{}, err
	}

	if !isValid {
		return BaseUrl{}, apierror.NewServiceBaseUrlError("invalid pattern base url")
	}

	isValid = baseUrlObj.isValidLength()

	if !isValid {
		return BaseUrl{}, apierror.NewServiceBaseUrlError("invalid length base url")
	}

	return baseUrlObj, nil
}

func (u BaseUrl) isValidPattern() (bool, error) {
	return regexp.MatchString(
		`^https?:\/\/[^\s\/$.?#].[^\s]*$`,
		u.value,
	)
}

func (u BaseUrl) isValidLength() bool {
	return len([]rune(u.value)) <= u.maxLength && len([]rune(u.value)) >= u.minLength
}

func (u BaseUrl) Value() string {
	return u.value
}
