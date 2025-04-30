package service

import (
	"context"
	"errors"
	"time"
	"github.com/yourusername/exchange-rate-service/internal/model"
	"github.com/yourusername/exchange-rate-service/internal/cache"
	"github.com/yourusername/exchange-rate-service/internal/utils"
)

// Service defines the business logic for exchange rate operations
type Service interface {
	Convert(ctx context.Context, req model.ConvertRequest) (model.ConvertResponse, error)
	GetRate(ctx context.Context, base, target, date string) (model.ExchangeRate, error)
	GetHistoricalRates(ctx context.Context, req model.HistoricalRatesRequest) (model.HistoricalRatesResponse, error)
}

// serviceImpl implements Service
// It uses a repository for external API, and an in-memory cache

type serviceImpl struct {
	repo  Repository
	cache *cache.RateCache
}

// Repository defines the methods for fetching rates from an external source
// This interface should be implemented in the repository package
type Repository interface {
	FetchRate(ctx context.Context, base, target, date string) (model.ExchangeRate, error)
}

// NewService constructs a Service
func NewService(repo Repository, cache *cache.RateCache) Service {
	return &serviceImpl{repo: repo, cache: cache}
}

// Convert performs currency conversion for a given amount and date
func (s *serviceImpl) Convert(ctx context.Context, req model.ConvertRequest) (model.ConvertResponse, error) {
	if !utils.IsValidCurrency(req.From, model.SupportedCurrencies) && !model.IsCryptoCurrency(req.From) {
		return model.ConvertResponse{}, errors.New("invalid 'from' currency")
	}
	if !utils.IsValidCurrency(req.To, model.SupportedCurrencies) && !model.IsCryptoCurrency(req.To) {
		return model.ConvertResponse{}, errors.New("invalid 'to' currency")
	}
	date := req.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if !utils.IsValidDate(date) {
		return model.ConvertResponse{}, errors.New("invalid date format")
	}
	if !utils.IsWithinLookback(date, 90) {
		return model.ConvertResponse{}, errors.New("date beyond 90-day lookback")
	}
	rate, err := s.GetRate(ctx, req.From, req.To, date)
	if err != nil {
		return model.ConvertResponse{}, err
	}
	return model.ConvertResponse{Amount: req.Amount * rate.Rate}, nil
}

// GetRate fetches the exchange rate for a currency pair and date
func (s *serviceImpl) GetRate(ctx context.Context, base, target, date string) (model.ExchangeRate, error) {
	// Try cache first
	if rate, ok := s.cache.Get(base, target, date); ok {
		return rate, nil
	}
	// Fetch from repository (external API)
	rate, err := s.repo.FetchRate(ctx, base, target, date)
	if err != nil {
		return model.ExchangeRate{}, err
	}
	s.cache.Set(rate)
	return rate, nil
}

// GetHistoricalRates fetches rates for a date range
func (s *serviceImpl) GetHistoricalRates(ctx context.Context, req model.HistoricalRatesRequest) (model.HistoricalRatesResponse, error) {
	if !utils.IsValidCurrency(req.From, model.SupportedCurrencies) && !model.IsCryptoCurrency(req.From) {
		return model.HistoricalRatesResponse{}, errors.New("invalid 'from' currency")
	}
	if !utils.IsValidCurrency(req.To, model.SupportedCurrencies) && !model.IsCryptoCurrency(req.To) {
		return model.HistoricalRatesResponse{}, errors.New("invalid 'to' currency")
	}
	if !utils.IsValidDate(req.StartDate) || !utils.IsValidDate(req.EndDate) {
		return model.HistoricalRatesResponse{}, errors.New("invalid date format")
	}
	if !utils.IsWithinLookback(req.StartDate, 90) || !utils.IsWithinLookback(req.EndDate, 90) {
		return model.HistoricalRatesResponse{}, errors.New("date beyond 90-day lookback")
	}
	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)
	if end.Before(start) {
		return model.HistoricalRatesResponse{}, errors.New("end date before start date")
	}
	var rates []model.ExchangeRate
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		rate, err := s.GetRate(ctx, req.From, req.To, dateStr)
		if err != nil {
			return model.HistoricalRatesResponse{}, err
		}
		rates = append(rates, rate)
	}
	return model.HistoricalRatesResponse{Rates: rates}, nil
}
