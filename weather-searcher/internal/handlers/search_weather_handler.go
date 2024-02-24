package handlers

import (
	"encoding/json"
	"net/http"

	usecases "github.com/ItaloG/go-weather-searcher/internal/usecases/weather"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type HttpSearchWeatherHandler struct {
	Usecase    usecases.SearchWeatherUseCase
	OTELTracer trace.Tracer
}

func NewHttpSearchWeatherHandler(searchWeatherUsecase usecases.SearchWeatherUseCase, tracer trace.Tracer) *HttpSearchWeatherHandler {
	return &HttpSearchWeatherHandler{
		Usecase:    searchWeatherUsecase,
		OTELTracer: tracer,
	}
}

func (h *HttpSearchWeatherHandler) Search(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.OTELTracer.Start(ctx, "search weather")
	defer span.End()

	cep := chi.URLParam(r, "cep")

	input := usecases.SearchWeatherInputDTO{
		Cep:      cep,
		OTELSpan: span,
		Ctx:      ctx,
	}

	output, err := h.Usecase.Execute(input)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
