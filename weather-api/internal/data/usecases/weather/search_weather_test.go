package usecases

import (
	"errors"
	"testing"

	services_protocols "github.com/ItaloG/go-weather-api/internal/data/protocols/services"
	services_mocks "github.com/ItaloG/go-weather-api/internal/data/usecases/mocks/services"
	validators_mocks "github.com/ItaloG/go-weather-api/internal/data/usecases/mocks/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldReturn_InvalidCepError_When_CepValidatorReturns_An_Error(t *testing.T) {
	cepValidatorMock := &validators_mocks.CepValidatorMock{}
	searchWeatherServiceMock := &services_mocks.SearchWeatherServiceMock{}

	cepValidatorMock.On("Validate", mock.Anything).Return(errors.New("zipcode error"))

	uc := NewSearchWeatherUseCase(cepValidatorMock, searchWeatherServiceMock)

	input := SearchWeatherInputDTO{Cep: "invalid_zipcode"}

	output, err := uc.Execute(input)

	assert.Nil(t, output)
	assert.Equal(t, err, ErrInvalidCep)
}

func TestShouldReturn_ZipCodeNotFoundError_When_SearchWeatherServiceReturns_An_Error(t *testing.T) {
	cepValidatorMock := &validators_mocks.CepValidatorMock{}
	searchWeatherServiceMock := &services_mocks.SearchWeatherServiceMock{}

	cepValidatorMock.On("Validate", mock.Anything).Return(nil)
	searchWeatherServiceMock.On("Search", mock.Anything, mock.Anything, mock.Anything).Return(&services_protocols.Weather{}, errors.New("weather error"))

	uc := NewSearchWeatherUseCase(cepValidatorMock, searchWeatherServiceMock)

	input := SearchWeatherInputDTO{Cep: "valid_zipcode"}

	output, err := uc.Execute(input)

	assert.Nil(t, output)
	assert.Equal(t, err, ErrZipCodeNotFound)
}

func TestShouldReturn_Weather_OnSuccess(t *testing.T) {
	cepValidatorMock := &validators_mocks.CepValidatorMock{}
	searchWeatherServiceMock := &services_mocks.SearchWeatherServiceMock{}

	cepValidatorMock.On("Validate", mock.Anything).Return(nil)

	var weatherMock *services_protocols.Weather = &services_protocols.Weather{City: "any_city", TempC: 100.0, TempF: 100.0, TempK: 100.0}
	searchWeatherServiceMock.On("Search", mock.Anything, mock.Anything, mock.Anything).Return(weatherMock, nil)

	uc := NewSearchWeatherUseCase(cepValidatorMock, searchWeatherServiceMock)

	input := SearchWeatherInputDTO{Cep: "valid_zipcode"}

	output, err := uc.Execute(input)

	assert.Equal(t, output.City, "any_city")
	assert.Equal(t, output.TempC, 100.0)
	assert.Equal(t, output.TempF, 100.0)
	assert.Equal(t, output.TempK, 100.0)
	assert.Nil(t, err)
}
