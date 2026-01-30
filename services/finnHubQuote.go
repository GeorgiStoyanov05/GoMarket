package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type finnhubQuoteResp struct {
	Current float64 `json:"c"`
}

func FetchCurrentPrice(symbol string) (float64, error) {
	token := os.Getenv("FINNHUB_API_KEY")
	if token == "" {
		return 0, errors.New("FINNHUB_API_KEY missing")
	}

	client := &http.Client{Timeout: 6 * time.Second}

	req, err := http.NewRequest("GET",
		"https://finnhub.io/api/v1/quote?symbol="+symbol+"&token="+token, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, errors.New("quote request failed")
	}

	var q finnhubQuoteResp
	if err := json.NewDecoder(resp.Body).Decode(&q); err != nil {
		return 0, err
	}
	if q.Current <= 0 {
		return 0, errors.New("invalid quote price")
	}

	return q.Current, nil
}
