package vo

import (
	"regexp"
	apierror "transport/api/src/error"
)

const (
	serviceBaseURLMaxLength = 255
	serviceBaseURLMinLength = 8
)

type ValidateLevel int

const (
	BaseURLValidateLevelStrict ValidateLevel = 1
	BaseURLValidateLevelNone   ValidateLevel = 0
)

type BaseURL struct {
	value                string
	maxLength, minLength int
}

func NewBaseURL(baseURL string, validateLevel ValidateLevel) (BaseURL, error) {
	baseURLObj := BaseURL{
		value:     baseURL,
		maxLength: serviceBaseURLMaxLength,
		minLength: serviceBaseURLMinLength,
	}

	if validateLevel == BaseURLValidateLevelStrict {
		isValid, err := baseURLObj.isValidPattern()

		if err != nil {
			return BaseURL{}, err
		}

		if !isValid {
			return BaseURL{}, apierror.NewServiceBaseURLError("invalid pattern base url")
		}

		isValid = baseURLObj.isValidLength()

		if !isValid {
			return BaseURL{}, apierror.NewServiceBaseURLError("invalid length base url")
		}
	}

	return baseURLObj, nil
}

func (u BaseURL) isValidPattern() (bool, error) {
	return regexp.MatchString(
		`^https?:\/\/[^\s\/$.?#].[^\s]*$`,
		u.value,
	)
}

func (u BaseURL) isValidLength() bool {
	return len([]rune(u.value)) <= u.maxLength && len([]rune(u.value)) >= u.minLength
}

func (u BaseURL) Value() string {
	return u.value
}
