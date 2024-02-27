package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	usecases "github.com/ItaloG/go-weather-api/internal/data/usecases/weather"
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
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.OTELTracer.Start(ctx, "input handler")
	defer span.End()

	var input usecases.SearchWeatherInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		span.SetAttributes(attribute.String("JSON Decode Error", err.Error()))
		span.SetAttributes(attribute.String("Response", fmt.Sprintf("Status: %d, Message: %s", http.StatusBadRequest, "you must send a zicode as string")))

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("you must send a zicode as string")
		return
	}

	span.SetAttributes(attribute.String("cep", input.Cep))

	input.OTELSpan = span
	input.Ctx = ctx

	output, err := h.Usecase.Execute(input)

	if err == usecases.ErrInvalidCep {
		span.SetAttributes(attribute.String("Response", fmt.Sprintf("Status: %d, Message: %s", http.StatusUnprocessableEntity, "invalid zipcode")))

		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode("invalid zipcode")
		return
	}
	if err == usecases.ErrZipCodeNotFound {
		span.SetAttributes(attribute.String("Response", fmt.Sprintf("Status: %d, Message: %s", http.StatusNotFound, "can not find zipcode")))

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("can not find zipcode")
		return
	}

	span.SetAttributes(attribute.String("UseCase Output", fmt.Sprintf("City: %s, TempC: %f, TempF: %f, TempK: %f", output.City, output.TempC, output.TempF, output.TempK)))

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		span.SetAttributes(attribute.String("JSON Encoder Error", err.Error()))
		span.SetAttributes(attribute.String("Response", fmt.Sprintf("Status: %d", http.StatusInternalServerError)))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
