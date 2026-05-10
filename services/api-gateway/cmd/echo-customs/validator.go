package echocustoms

import (
	sharedValidator "github.com/dinno7/ride-sharing/shared/validator"
)

type EchoValidator struct{}

func NewEchoValidator() *EchoValidator {
	return &EchoValidator{}
}

func (cv *EchoValidator) Validate(i any) error {
	return sharedValidator.ValidateData(i)
}
