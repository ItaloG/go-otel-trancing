package validators

import (
	"errors"
	"regexp"
)

type CepValidator struct {
}

func NewCepValidator() *CepValidator {
	return &CepValidator{}
}

func (v *CepValidator) Validate(cep string) error {

	if cep == "" {
		return errors.New("empty zipcode")
	}
	isValid, err := regexp.MatchString(`^\d{8}$`, cep)
	if err != nil || !isValid {
		return errors.New("invalid zipcode format")
	}

	return nil
}
