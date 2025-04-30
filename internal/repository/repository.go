package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/exchange-rate-service/config"
	"github.com/yourusername/exchange-rate-service/internal/model"
)

type APIRepository struct {
	BaseURL string
	Client  *http.Client
	APIKey  string
}

func NewAPIRepository() *APIRepository {
	config.InitConfig()
	return &APIRepository{
		BaseURL: "https://api.exchangerate.host",
		Client:  &http.Client{Timeout: 10 * time.Second},
		APIKey:  config.Config.ExchangeAPIKey,
	}
}

// FetchRate fetches the exchange rate for a currency pair and date from exchangerate.host
func (r *APIRepository) FetchRate(ctx context.Context, base, target, date string) (model.ExchangeRate, error) {
	if r.APIKey == "" {
		return model.ExchangeRate{}, errors.New("EXCHANGE_API_KEY not set in environment")
	}
	var url string
	if date == time.Now().Format("2006-01-02") {
		url = fmt.Sprintf("%s/live?access_key=%s&source=%s&currencies=%s", r.BaseURL, r.APIKey, base, target)
	} else {
		url = fmt.Sprintf("%s/historical?access_key=%s&date=%s&source=%s&currencies=%s", r.BaseURL, r.APIKey, date, base, target)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return model.ExchangeRate{}, err
	}
	resp, err := r.Client.Do(req)
	if err != nil {
		return model.ExchangeRate{}, err
	}
	defer resp.Body.Close()
	var apiResp struct {
		Success bool `json:"success"`
		Error   *struct {
			Code int    `json:"code"`
			Info string `json:"info"`
		} `json:"error,omitempty"`
		Source string             `json:"source"`
		Quotes map[string]float64 `json:"quotes"`
		Date   string             `json:"date"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return model.ExchangeRate{}, err
	}
	if !apiResp.Success {
		errMsg := "API error"
		if apiResp.Error != nil {
			errMsg = fmt.Sprintf("API error %d: %s", apiResp.Error.Code, apiResp.Error.Info)
		}
		return model.ExchangeRate{}, errors.New(errMsg)
	}
	pair := base + target
	rate, ok := apiResp.Quotes[pair]
	if !ok {
		return model.ExchangeRate{}, errors.New("target currency not found in API response")
	}
	return model.ExchangeRate{
		Base:      base,
		Target:    target,
		Rate:      rate,
		Date:      apiResp.Date,
		Source:    "exchangerate.host",
		UpdatedAt: time.Now(),
	}, nil
}
