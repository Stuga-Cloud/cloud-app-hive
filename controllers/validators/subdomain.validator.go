package validators

import (
	"github.com/go-playground/validator"
	"regexp"
)

func IsACustomStringForSubdomainValidation(fl validator.FieldLevel) bool {
	// Since Namespace and name will be in a url, they should not contain special characters like / or . only a-z, A-Z, 0-9, _ and -
	match, err := regexp.MatchString("^[a-zA-Z0-9_-]*$", fl.Field().String())
	if err != nil {
		return false
	}
	return match
}
