package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	services_protocols "github.com/ItaloG/go-weather-searcher/internal/protocols/services"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"github.com/rainycape/unidecode"
)

type WeatherApiResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree float64 `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   float64 `json:"humidity"`
		Cloud      float64 `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

type SearchWeatherService struct {
	ClientUrl    string
	WeatherToken string
	OTELTracer   trace.Tracer
}

func NewSearchWeatherService(clientUrl, weatherToken string, tracer trace.Tracer) *SearchWeatherService {
	return &SearchWeatherService{ClientUrl: clientUrl, WeatherToken: weatherToken, OTELTracer: tracer}
}

var ErrWeatherNotFound = errors.New("can not found weather")

func (sw *SearchWeatherService) Search(ctx context.Context, location string, OTELSpan trace.Span) (*services_protocols.Weather, error) {
	ctx, span := sw.OTELTracer.Start(ctx, "request to weather API")
	defer span.End()

	startTime := time.Now()

	formattedLocation := strings.Replace(location, " ", "_", -1)
	formattedLocation = unidecode.Unidecode(formattedLocation)

	OTELSpan.SetAttributes(attribute.String("HTTP Request", sw.ClientUrl+formattedLocation+"&key="+sw.WeatherToken))

	req, err := http.NewRequestWithContext(ctx, "GET", sw.ClientUrl+formattedLocation+"&key="+sw.WeatherToken, nil)
	if err != nil {
		OTELSpan.SetAttributes(attribute.String("HTTP Error", err.Error()))
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		OTELSpan.SetAttributes(attribute.String("HTTP Do Error", err.Error()))
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, ErrWeatherNotFound
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		OTELSpan.SetAttributes(attribute.String("HTTP Body Error", err.Error()))
		return nil, err
	}

	var weatherResponse WeatherApiResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		OTELSpan.SetAttributes(attribute.String("JSON Unmarshal Error", err.Error()))
		return nil, err
	}

	tempK := float64(weatherResponse.Current.TempC) + 273

	requestDuration := time.Since(startTime)
	span.SetAttributes(attribute.Float64("weather API request duration", requestDuration.Seconds()))

	weatherapiResponse, err:= json.Marshal(weatherResponse)
	if err != nil {
		return nil, err
	}
	OTELSpan.SetAttributes(attribute.String("HTTP Responses", string(weatherapiResponse)))

	return &services_protocols.Weather{
		TempC: float64(weatherResponse.Current.TempC),
		TempF: float64(weatherResponse.Current.TempF),
		TempK: tempK,
	}, nil
}
