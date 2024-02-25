package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	usecases "github.com/ItaloG/go-weather-searcher/internal/usecases/weather"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	cep := chi.URLParam(r, "cep")

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.OTELTracer.Start(ctx, "search weather")
	defer span.End()

	input := usecases.SearchWeatherInputDTO{
		Cep:      cep,
		OTELSpan: span,
		Ctx:      ctx,
	}

	output, err := h.Usecase.Execute(input)
	if err == usecases.ErrInvalidCep {
		span.SetAttributes(attribute.String("Response", fmt.Sprintf("Status: %d, Message: %s", http.StatusUnprocessableEntity, "invalid zipcode")))

		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode("invalid zipcode")
		return
	}
	if err == usecases.ErrLocationNotFound {
		span.SetAttributes(attribute.String("Response", fmt.Sprintf("Status: %d, Message: %s", http.StatusNotFound, "can not find zipcode")))

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("can not find zipcode")
		return
	}
	if err == usecases.ErrWeatherNotFound {
		span.SetAttributes(attribute.String("Response", fmt.Sprintf("Status: %d, Message: %s", http.StatusNotFound, "can not find weather")))

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("can not find weather")
		return
	}

	span.SetAttributes(attribute.String("UseCase Output", fmt.Sprintf("City: %s, TempC: %f, TempF: %f, TempK: %f", output.City, output.TempC, output.TempF, output.TempK)))

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		span.SetAttributes(attribute.String("JSON Encoder Error", err.Error()))
		span.SetAttributes(attribute.String("Response", fmt.Sprintf("Status: %d", http.StatusInternalServerError)))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
