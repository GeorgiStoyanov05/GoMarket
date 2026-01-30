package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type FinnhubSearchItem struct {
	Description   string `json:"description"`
	DisplaySymbol string `json:"displaySymbol"`
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
}

type finnhubSearchResp struct {
	Count  int                `json:"count"`
	Result []FinnhubSearchItem `json:"result"`
}

func SearchSymbols(query string, limit int) ([]FinnhubSearchItem, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return []FinnhubSearchItem{}, nil
	}

	token := os.Getenv("FINNHUB_API_KEY")
	if token == "" {
		return nil, errors.New("FINNHUB_API_KEY missing")
	}

	if limit <= 0 {
		limit = 10
	}

	client := &http.Client{Timeout: 6 * time.Second}

	u := "https://finnhub.io/api/v1/search?q=" + url.QueryEscape(q) + "&token=" + url.QueryEscape(token)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("finnhub search failed: status %d", resp.StatusCode)
	}

	var out finnhubSearchResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	results := out.Result
	// basic cleanup: drop empty symbols
	clean := results[:0]
	for _, r := range results {
		if strings.TrimSpace(r.Symbol) != "" {
			clean = append(clean, r)
		}
	}
	results = clean

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}
