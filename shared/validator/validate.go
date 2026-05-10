package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var DefaultValidator *validator.Validate

func init() {
	DefaultValidator = validator.New(validator.WithRequiredStructEnabled())
}

func ValidateHttpBody[T any](req *http.Request) (*T, error) {
	data, err := readJSON[T](req)
	if err != nil {
		return nil, errors.New("faild to read request's body")
	}

	if err := ValidateData(data); err != nil {
		return nil, err
	}

	return data, nil
}

func ValidateData[T any](data T) error {
	if err := DefaultValidator.Struct(data); err != nil {
		msg := "failed validate fields"
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return err
		}

		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			errorStrs := []string{}
			for _, e := range validateErrs {
				errorMsg := fmt.Sprintf("%s: %s", e.Field(), e.Tag())
				if e.Param() != "" {
					errorMsg += "=" + e.Param()
				}
				errorStrs = append(errorStrs, errorMsg)
			}
			msg += fmt.Sprintf(": %s", strings.Join(errorStrs, ", "))
		}
		return errors.New(msg)
	}

	return nil
}

func readJSON[T any](req *http.Request) (*T, error) {
	var data T
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
