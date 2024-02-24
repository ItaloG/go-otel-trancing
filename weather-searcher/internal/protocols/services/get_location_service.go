package services_protocols

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Location struct {
	Localidade string `json:"localidade"`
}

type GetLocationService interface {
	GetLocation(ctx context.Context, cep string, OTELSpan trace.Span) (*Location, error)
}
