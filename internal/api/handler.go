package api

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/yourusername/exchange-rate-service/internal/model"
	"github.com/yourusername/exchange-rate-service/internal/service"
)

// Endpoints struct for go-kit

type Endpoints struct {
	ConvertEndpoint          endpoint.Endpoint
	GetRateEndpoint          endpoint.Endpoint
	GetHistoricalRatesEndpoint endpoint.Endpoint
}

// MakeEndpoints creates go-kit endpoints from the service
func MakeEndpoints(svc service.Service) Endpoints {
	return Endpoints{
		ConvertEndpoint:          makeConvertEndpoint(svc),
		GetRateEndpoint:          makeGetRateEndpoint(svc),
		GetHistoricalRatesEndpoint: makeGetHistoricalRatesEndpoint(svc),
	}
}

func makeConvertEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.ConvertRequest)
		return svc.Convert(ctx, req)
	}
}

func makeGetRateEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.ConvertRequest)
		rate, err := svc.GetRate(ctx, req.From, req.To, req.Date)
		if err != nil {
			return nil, err
		}
		return rate, nil
	}
}

func makeGetHistoricalRatesEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.HistoricalRatesRequest)
		return svc.GetHistoricalRates(ctx, req)
	}
}

// HTTP Handlers

func MakeHTTPHandler(endpoints Endpoints) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/convert", kithttp.NewServer(
		endpoints.ConvertEndpoint,
		decodeConvertRequest,
		encodeResponse,
	))
	mux.Handle("/rate", kithttp.NewServer(
		endpoints.GetRateEndpoint,
		decodeConvertRequest,
		encodeResponse,
	))
	mux.Handle("/history", kithttp.NewServer(
		endpoints.GetHistoricalRatesEndpoint,
		decodeHistoricalRatesRequest,
		encodeResponse,
	))
	return mux
}

func decodeConvertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amount := 1.0
	if a := r.URL.Query().Get("amount"); a != "" {
		// Optionally parse amount
	}
	date := r.URL.Query().Get("date")
	return model.ConvertRequest{From: from, To: to, Amount: amount, Date: date}, nil
}

func decodeHistoricalRatesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	return model.HistoricalRatesRequest{From: from, To: to, StartDate: start, EndDate: end}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
