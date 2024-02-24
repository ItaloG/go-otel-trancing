package services_mocks

import (
	"context"

	services_protocols "github.com/ItaloG/go-weather-api/internal/data/protocols/services"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace"
)

type SearchWeatherServiceMock struct {
	mock.Mock
}

func (m *SearchWeatherServiceMock) Search(ctx context.Context, cep string, OTELSpan trace.Span) (*services_protocols.Weather, error) {
	args := m.Called(ctx, cep, OTELSpan)
	return args.Get(0).(*services_protocols.Weather), args.Error(1)
}
