package main

import (
	"log"
	"net/http"

	"github.com/ItaloG/go-weather-searcher/internal/handlers"
	"github.com/ItaloG/go-weather-searcher/internal/infra/services"
	usecases "github.com/ItaloG/go-weather-searcher/internal/usecases/weather"
	"github.com/ItaloG/go-weather-searcher/internal/validators"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func init() {
	viper.AutomaticEnv()
}

func main() {
	initOpenTelemetry(viper.GetString("ZIPKIN_URL"), viper.GetString("SERVICE_NAME"))

	tracer := otel.Tracer(viper.GetString("SERVICE_NAME"))

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(otelhttp.NewMiddleware("OTEL HTTP Middleware" + viper.GetString("SERVICE_NAME")))

	searchWeatherService := services.NewSearchWeatherService(viper.GetString("WEATHER_API"), viper.GetString("WEATHER_TOKEN"))
	getLocationService := services.NewGetLocationService(viper.GetString("LOCATION_API"))
	cepValidator := validators.NewCepValidator()
	searchWeatherUsecase := usecases.NewSearchWeatherUseCase(cepValidator, getLocationService, searchWeatherService)
	searchWeatherHandler := handlers.NewHttpSearchWeatherHandler(*searchWeatherUsecase, tracer)

	r.Get("/{cep}", searchWeatherHandler.Search)

	http.ListenAndServe(viper.GetString("SERVER_PORT"), r)
}

func initOpenTelemetry(zipkinUrl, serviceName string) {
	exporter, err := zipkin.New(zipkinUrl)
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	))
}
