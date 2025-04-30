package repository

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/exchange-rate-service/internal/model"
)

type mockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestFetchRate_Success(t *testing.T) {
	jsonResp := `{"success":true,"source":"USD","quotes":{"USDINR":83.0},"date":"2025-04-30"}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jsonResp))
	}))
	defer ts.Close()

	repo := &APIRepository{
		BaseURL: ts.URL,
		Client:  ts.Client(),
		APIKey:  "dummy",
	}
	ctx := context.Background()
	rate, err := repo.FetchRate(ctx, "USD", "INR", "2025-04-30")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if rate.Rate != 83.0 {
		t.Errorf("Expected rate 83.0, got %v", rate.Rate)
	}
	if rate.Base != "USD" || rate.Target != "INR" {
		t.Errorf("Base/Target mismatch: got %s/%s", rate.Base, rate.Target)
	}
}

func TestFetchRate_APIFailure(t *testing.T) {
	jsonResp := `{"success":false,"error":{"code":101,"info":"Invalid key"}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jsonResp))
	}))
	defer ts.Close()

	repo := &APIRepository{
		BaseURL: ts.URL,
		Client:  ts.Client(),
		APIKey:  "badkey",
	}
	ctx := context.Background()
	_, err := repo.FetchRate(ctx, "USD", "INR", "2025-04-30")
	if err == nil || !errors.Is(err, err) {
		t.Errorf("Expected error for API failure, got %v", err)
	}
}

func TestFetchRate_MissingKey(t *testing.T) {
	repo := &APIRepository{
		BaseURL: "http://example.com",
		Client:  &http.Client{Timeout: 1 * time.Second},
		APIKey:  "",
	}
	ctx := context.Background()
	_, err := repo.FetchRate(ctx, "USD", "INR", "2025-04-30")
	if err == nil {
		t.Error("Expected error for missing API key")
	}
}