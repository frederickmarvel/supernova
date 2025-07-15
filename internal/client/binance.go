package client

import (
	"encoding/json"
	"net/http"
	"time"
)

type Kline struct {
	OpenTime  int64
	Open      string
	High      string
	Low       string
	Close     string
	Volume    string
	CloseTime int64
}

// Fetch last 180 days of daily klines
func FetchKlines(symbol string) ([]Kline, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest(
		http.MethodGet,
		"https://api.binance.com/api/v3/klines",
		nil,
	)
	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("interval", "1d")
	q.Add("limit", "180")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	klines := make([]Kline, len(raw))
	for i, entry := range raw {
		klines[i] = Kline{
			OpenTime:  int64(entry[0].(float64)),
			Open:      entry[1].(string),
			High:      entry[2].(string),
			Low:       entry[3].(string),
			Close:     entry[4].(string),
			Volume:    entry[5].(string),
			CloseTime: int64(entry[6].(float64)),
		}
	}
	return klines, nil
}
