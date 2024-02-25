package main

import (
	"log"
	"net/http"

	usecases "github.com/ItaloG/go-weather-api/internal/data/usecases/weather"
	"github.com/ItaloG/go-weather-api/internal/http/handlers"
	"github.com/ItaloG/go-weather-api/internal/infra/services"
	"github.com/ItaloG/go-weather-api/internal/validators"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
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

	searchWeatherService := services.NewSearchWeatherService(viper.GetString("WEATHER_API"))
	// searchWeatherService := services.NewSearchWeatherService("http://172.23.128.46:8081/")
	cepValidator := validators.NewCepValidator()
	searchWeatherUsecase := usecases.NewSearchWeatherUseCase(cepValidator, searchWeatherService)
	searchWaterHandler := handlers.NewHttpSearchWeatherHandler(*searchWeatherUsecase, tracer)

	r.Post("/weather", searchWaterHandler.Search)

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
