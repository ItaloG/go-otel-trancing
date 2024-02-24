package services_protocols

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Weather struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type SearchWeatherService interface {
	Search(ctx context.Context, cep string, OTELSpan trace.Span) (*Weather, error)
}
