package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"math"
	"net/http"
)

type InputCEP struct {
	Cep string `json:"cep"`
}
type ApiViaCepResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

type GetWeatherResponse struct {
	Celsius    float64 `json:"celsius"`
	Fahrenheit float64 `json:"fahrenheit"`
	Kelvin     float64 `json:"kelvin"`
}

type ApiWeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

var tracer trace.Tracer

func main() {
	initTracer()

	mux := http.NewServeMux()
	mux.HandleFunc("/service-B", handlerServiceB)
	err := http.ListenAndServe(":8091", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func initTracer() {

	exporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider( //  Cria um trace provider para gerenciar todos os spans dentro do serviço
		sdktrace.WithBatcher(exporter), // Defini o exportador que enviara os spans para o zipkin
		sdktrace.WithResource(resource.Default()),
	)
	otel.SetTracerProvider(tp) // defini o tracer provider global
	otel.SetTextMapPropagator(propagation.TraceContext{})
	tracer = tp.Tracer("service-B") // Cria o tracer para esse servico
}

func handlerServiceB(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header)) // extraiu o ctx
	ctx, span := tracer.Start(ctx, "handlerServiceB")
	defer span.End()
	var input InputCEP
	err := json.NewDecoder(r.Body).Decode(&input.Cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	location, err := getLocation(input.Cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	celsiuTemp, err := getCurrentCelsiusTemp(ctx, location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fahrenheintTemp := celsiuTemp*1.8 + 32
	kelvin := celsiuTemp + 273

	output := GetWeatherResponse{
		roundToTwo(celsiuTemp),
		roundToTwo(fahrenheintTemp),
		roundToTwo(kelvin),
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)

}

func getCurrentCelsiusTemp(ctx context.Context, location string) (float64, error) {
	ctx, span := tracer.Start(ctx, "getCurrentCelsiusTemp")
	defer span.End()
	apiKey := "f875c284c1114aec9c5220427250402"

	url_get_weather := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, location)

	resp, err := http.Get(url_get_weather)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var apiWeatherResponse ApiWeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&apiWeatherResponse)
	if err != nil {
		return 0, err
	}
	return apiWeatherResponse.Current.TempC, nil
}
func roundToTwo(value float64) float64 {
	return math.Round(value*100) / 100
}

func getLocation(cep string) (string, error) {
	url_get_loc := "https://viacep.com.br/ws/" + cep + "/json"

	res, err := http.Get(url_get_loc)
	if err != nil {

		return "", errors.New("Error getting location")
	}
	defer res.Body.Close()
	var output ApiViaCepResponse

	if err = json.NewDecoder(res.Body).Decode(&output); err != nil {
		return "", errors.New("Error getting location")
	}
	return output.Localidade, nil
}
