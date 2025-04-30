package model

import "time"

// Supported fiat and crypto currencies
var SupportedCurrencies = []string{"USD", "INR", "EUR", "JPY", "GBP"}
var SupportedCryptoCurrencies = []string{"BTC", "ETH", "USDT"}

// IsCryptoCurrency returns true if the code is a supported crypto
func IsCryptoCurrency(code string) bool {
	for _, c := range SupportedCryptoCurrencies {
		if c == code {
			return true
		}
	}
	return false
}

// ExchangeRate represents the exchange rate between two currencies on a specific date
// Rate is always with respect to 1 unit of base currency
// Example: 1 USD = 83.12 INR
// Rate = 83.12, Base = USD, Target = INR
// Date is always in YYYY-MM-DD format
// Source is the data provider (e.g., exchangerate.host)
type ExchangeRate struct {
	Base      string    `json:"base"`
	Target    string    `json:"target"`
	Rate      float64   `json:"rate"`
	Date      string    `json:"date"`
	Source    string    `json:"source"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ConvertRequest represents a currency conversion request
type ConvertRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date,omitempty"`
}

// ConvertResponse represents a currency conversion response
type ConvertResponse struct {
	Amount float64 `json:"amount"`
}

// HistoricalRatesRequest represents a request for historical rates
type HistoricalRatesRequest struct {
	From      string `json:"from"`
	To        string `json:"to"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// HistoricalRatesResponse represents a response for historical rates
type HistoricalRatesResponse struct {
	Rates []ExchangeRate `json:"rates"`
}
