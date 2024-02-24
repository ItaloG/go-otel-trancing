package validators_mocks

import "github.com/stretchr/testify/mock"

type CepValidatorMock struct {
	mock.Mock
}

func (m *CepValidatorMock) Validate(cep string) error {
	args := m.Called(cep)
	return args.Error(0)
}
