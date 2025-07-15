package service

import (
	"database/sql"
	"math"
	"strconv"
	"time"

	"github.com/frederickmarvel/supernova/internal/client"
)

var lambdas = []float64{0.5, 0.757858283, 0.870550563, 0.933032992, 0.965936329, 0.982820599}
var nfs = []float64{1.0000, 1.0000, 1.0000, 1.0000, 1.0020, 1.0462}

func ewma(prices []float64, lam, nf float64) float64 {
	weights := make([]float64, len(prices))
	var sumW float64
	for i := range prices {
		w := (1 - lam) * math.Pow(lam, float64(i))
		weights[i] = w
		sumW += w
	}
	norm := 1 / sumW
	var total float64
	for i, price := range prices {
		total += weights[i] * price * nf
	}
	return norm * total
}

func computeIndicator(symbol string) (float64, error) {
	klines, err := client.FetchKlines(symbol)
	if err != nil {
		return 0, err
	}
	if len(klines) < 180 {
		return 0, nil
	}
	closes := make([]float64, len(klines))
	for i := range klines {
		c, _ := strconv.ParseFloat(klines[len(klines)-1-i].Close, 64)
		closes[i] = c
	}

	mas := []float64{
		ewma(closes, lambdas[0], nfs[0]),
		ewma(closes, lambdas[1], nfs[1]),
		ewma(closes, lambdas[2], nfs[2]),
		ewma(closes, lambdas[3], nfs[3]),
		ewma(closes, lambdas[4], nfs[4]),
		ewma(closes, lambdas[5], nfs[5]),
	}
	diffs := []float64{
		mas[0] - mas[2],
		mas[1] - mas[3],
		mas[2] - mas[4],
		mas[3] - mas[5],
	}
	var sumSigns float64
	for _, d := range diffs {
		if d >= 0 {
			sumSigns += 1
		} else {
			sumSigns -= 1
		}
	}
	return sumSigns / 4, nil
}

func UpdateTrends(db *sql.DB) error {
	now := time.Now()
	symbols := map[string]string{
		"bitcoin":  "BTCUSDT",
		"ethereum": "ETHUSDT",
		"solana":   "SOLUSDT",
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		`CREATE TABLE IF NOT EXISTS trend_indicator (
            bitcoin_trend FLOAT,
            ethereum_trend FLOAT,
            solana_trend FLOAT,
            timestamp TIMESTAMP
        )`,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	var vals []interface{}
	for _, sym := range symbols {
		ind, err := computeIndicator(sym)
		if err != nil {
			tx.Rollback()
			return err
		}
		vals = append(vals, ind)
	}
	vals = append(vals, now)
	_, err = tx.Exec(
		`INSERT INTO trend_indicator (
            bitcoin_trend, ethereum_trend, solana_trend, timestamp
        ) VALUES ($1,$2,$3,$4)`,
		vals...,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func GetLatest(db *sql.DB) (map[string]float64, time.Time, error) {
	row := db.QueryRow(
		`SELECT bitcoin_trend, ethereum_trend, solana_trend, timestamp
         FROM trend_indicator ORDER BY timestamp DESC LIMIT 1`,
	)
	var b, e, s float64
	var ts time.Time
	err := row.Scan(&b, &e, &s, &ts)
	if err != nil {
		return nil, time.Time{}, err
	}
	return map[string]float64{
		"bitcoin":  b,
		"ethereum": e,
		"solana":   s,
	}, ts, nil
}
