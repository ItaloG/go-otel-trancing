package usecases

import (
	"context"
	"errors"

	services_protocols "github.com/ItaloG/go-weather-searcher/internal/protocols/services"
	validators_protocols "github.com/ItaloG/go-weather-searcher/internal/protocols/validators"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SearchWeatherInputDTO struct {
	Cep      string `json:"cep"`
	OTELSpan trace.Span
	Ctx      context.Context
}

type SearchWeatherOutputDTO struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type SearchWeatherUseCase struct {
	CepValidator         validators_protocols.CepValidator
	GetLocationService   services_protocols.GetLocationService
	SearchWeatherService services_protocols.SearchWeatherService
}

func NewSearchWeatherUseCase(cepValidator validators_protocols.CepValidator, getLocationService services_protocols.GetLocationService, searchWeatherService services_protocols.SearchWeatherService) *SearchWeatherUseCase {
	return &SearchWeatherUseCase{
		CepValidator:         cepValidator,
		GetLocationService:   getLocationService,
		SearchWeatherService: searchWeatherService,
	}
}

var ErrInvalidCep = errors.New("invalid zipcode")
var ErrLocationNotFound = errors.New("can not find zipcode")
var ErrWeatherNotFound = errors.New("can not found weather")

func (uc *SearchWeatherUseCase) Execute(input SearchWeatherInputDTO) (*SearchWeatherOutputDTO, error) {
	err := uc.CepValidator.Validate(input.Cep)
	if err != nil {
		input.OTELSpan.SetAttributes(attribute.String("Validate Cep Error", err.Error()))
		return nil, ErrInvalidCep
	}

	location, err := uc.GetLocationService.GetLocation(input.Ctx, input.Cep, input.OTELSpan)
	if err != nil {
		return nil, ErrLocationNotFound
	}

	weather, err := uc.SearchWeatherService.Search(input.Ctx, location.Localidade, input.OTELSpan)
	if err != nil {
		return nil, ErrWeatherNotFound
	}

	return &SearchWeatherOutputDTO{
		City:  location.Localidade,
		TempC: weather.TempC,
		TempF: weather.TempF,
		TempK: weather.TempK,
	}, nil
}
