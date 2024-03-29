package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	services_protocols "github.com/ItaloG/go-weather-api/internal/data/protocols/services"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type SearchWeatherService struct {
	ClientUrl  string
	OTELTracer trace.Tracer
}

func NewSearchWeatherService(clientUrl string, tracer trace.Tracer) *SearchWeatherService {
	return &SearchWeatherService{ClientUrl: clientUrl, OTELTracer: tracer}
}

var ErrWeatherNotFound = errors.New("can not found weather")

func (sw *SearchWeatherService) Search(ctx context.Context, cep string, OTELSpan trace.Span) (*services_protocols.Weather, error) {
	ctx, span := sw.OTELTracer.Start(ctx, "request to weather searcher microservice")
	defer span.End()

	startTime := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", sw.ClientUrl+cep, nil)
	if err != nil {
		OTELSpan.SetAttributes(attribute.String("HTTP Error", err.Error()))
		return &services_protocols.Weather{}, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		OTELSpan.SetAttributes(attribute.String("HTTP Do Error", err.Error()))
		return &services_protocols.Weather{}, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &services_protocols.Weather{}, ErrWeatherNotFound
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		OTELSpan.SetAttributes(attribute.String("HTTP Body Error", err.Error()))
		return &services_protocols.Weather{}, err
	}
	var weather *services_protocols.Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		OTELSpan.SetAttributes(attribute.String("JSON Unmarshal Error", err.Error()))
		return &services_protocols.Weather{}, err
	}

	requestDuration := time.Since(startTime)
	span.SetAttributes(attribute.Float64("weather searcher request duration", requestDuration.Seconds()))

	return weather, nil
}
