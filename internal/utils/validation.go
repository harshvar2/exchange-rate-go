package utils

import (
	"regexp"
	"strings"
	"time"
)

// TODO: Implement validation helpers for currency codes, dates, etc.

var currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)

// IsValidCurrency checks if the currency code is valid and supported
func IsValidCurrency(code string, supported []string) bool {
	code = strings.ToUpper(code)
	if !currencyPattern.MatchString(code) {
		return false
	}
	for _, c := range supported {
		if c == code {
			return true
		}
	}
	return false
}

// IsValidDate checks if the date is in YYYY-MM-DD format and logical
func IsValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// IsWithinLookback checks if the date is within the allowed lookback window (e.g., 90 days)
func IsWithinLookback(date string, days int) bool {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	limit := time.Now().AddDate(0, 0, -days)
	return !parsed.Before(limit)
}
