package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

type InputCEP struct {
	Cep string `json:"cep"`
}
type GetWeatherResponse struct {
	Celsius    float64 `json:"celsius"`
	Fahrenheit float64 `json:"fahrenheit"`
	Kelvin     float64 `json:"kelvin"`
}

var tracer trace.Tracer

func initTracer() {
	exporter, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider( //  Cria um trace provider para gerenciar todos os spans dentro do serviço
		sdktrace.WithBatcher(exporter), // Defini o exportador que enviara os spans para o zipkin
		sdktrace.WithResource(resource.Default()),
	)
	otel.SetTracerProvider(tp) // defini o tracer provider global
	otel.SetTextMapPropagator(propagation.TraceContext{})
	tracer = tp.Tracer("service-A") // cria o tracer para esse servico
}

func main() {
	initTracer()

	mux := http.NewServeMux()
	mux.HandleFunc("/service-A", HandlerServiceA)

	wrappedMux := otelhttp.NewHandler(mux, "service-A") // middleare para que toda req http seja automaticamente rastreada(handler instrumentado)
	if err := http.ListenAndServe(":8010", wrappedMux); err != nil {
		log.Fatal(err)
	}

}

func HandlerServiceA(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "HandlerServiceA")
	defer span.End()
	time.Sleep(2 * time.Second)
	var inputCEP InputCEP
	err := json.NewDecoder(r.Body).Decode(&inputCEP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = ValidateCEP(inputCEP.Cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(inputCEP.Cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", "http://service-b:8091/service-B", bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var responseServiceB GetWeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&responseServiceB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseServiceB)

}

func ValidateCEP(cep string) error {
	if len(cep) != 8 {
		return errors.New("invalid cep")
	}
	return nil
}
