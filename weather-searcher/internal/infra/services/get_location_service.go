package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	services_protocols "github.com/ItaloG/go-weather-searcher/internal/protocols/services"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        bool   `json:"erro"`
}

type GetLocationService struct {
	ClientUrl string
}

func NewGetLocationService(clientUrl string) *GetLocationService {
	return &GetLocationService{ClientUrl: clientUrl}
}

var ErrLocationNotFound = errors.New("can not found location")

func (s *GetLocationService) GetLocation(ctx context.Context, cep string, OTELSpan trace.Span) (*services_protocols.Location, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.ClientUrl+cep+"/json/", nil)
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
	if resp.StatusCode != 200 {
		return nil, ErrLocationNotFound
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		OTELSpan.SetAttributes(attribute.String("HTTP Body Error", err.Error()))
		return nil, err
	}
	var c ViaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		OTELSpan.SetAttributes(attribute.String("JSON Unmarshal Error", err.Error()))
		return nil, err
	}
	if c.Erro {
		OTELSpan.SetAttributes(attribute.String("Vicep Error", "Cep inválido"))
		return nil, ErrLocationNotFound
	}

	return &services_protocols.Location{Localidade: c.Localidade}, nil
}