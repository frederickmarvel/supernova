package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

const (
	apiURL = "https://indodax.com/tapi"
)

type Config struct {
	APIKey    string
	APISecret string
}

type IndodaxClient struct {
	config     Config
	httpClient *http.Client
	nonce      int64
	mu         sync.Mutex
}

func NewClient(apiKey, apiSecret string, client *http.Client) *IndodaxClient {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &IndodaxClient{
		config:     Config{APIKey: apiKey, APISecret: apiSecret},
		httpClient: client,
		nonce:      time.Now().UnixNano() / int64(time.Millisecond),
	}
}

func (c *IndodaxClient) nextNonce() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nonce++
	return c.nonce
}

func (c *IndodaxClient) doSigned(ctx context.Context, method string, params map[string]string, out interface{}) error {
	if params == nil {
		params = make(map[string]string)
	}
	params["method"] = method
	params["nonce"] = strconv.FormatInt(c.nextNonce(), 10)

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}
	body := form.Encode()

	mac := hmac.New(sha512.New, []byte(c.config.APISecret))
	mac.Write([]byte(body))
	sig := hex.EncodeToString(mac.Sum(nil))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Key", c.config.APIKey)
	req.Header.Set("Sign", sig)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(500 * time.Millisecond)
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("retry request: %w", err)
		}
		defer resp.Body.Close()
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http %d: %s", resp.StatusCode, data)
	}

	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	val := reflect.ValueOf(out).Elem()
	f := val.FieldByName("Success")
	if f.IsValid() && f.Int() == 0 {
		errFld := val.FieldByName("Error")
		return fmt.Errorf("api error: %s", errFld.String())
	}

	return nil
}

func (c *IndodaxClient) GetInfo(ctx context.Context) (*GetInfoResponse, error) {
	var res GetInfoResponse
	if err := c.doSigned(ctx, "getInfo", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) GetTransactionHistory(ctx context.Context, start, end string) (*TransHistoryResponse, error) {
	params := map[string]string{}
	if start != "" {
		params["start"] = start
	}
	if end != "" {
		params["end"] = end
	}
	var res TransHistoryResponse
	if err := c.doSigned(ctx, "transHistory", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) Trade(ctx context.Context, pair, tradeType string, price, amount, idrAmount decimal.Decimal, clientID, tif string) (*TradeResponse, error) {
	params := map[string]string{"pair": pair, "type": tradeType}
	if !price.IsZero() {
		params["price"] = price.String()
	}
	if !amount.IsZero() {
		params["btc"] = amount.String()
	}
	if !idrAmount.IsZero() {
		params["idr"] = idrAmount.String()
	}
	if clientID != "" {
		params["client_order_id"] = clientID
	}
	if tif != "" {
		params["time_in_force"] = tif
	}
	var res TradeResponse
	if err := c.doSigned(ctx, "trade", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) GetTradeHistory(ctx context.Context, pair string, opts ...TradeHistoryOption) (*TradeHistoryResponse, error) {
	params := map[string]string{"pair": pair}
	for _, opt := range opts {
		opt(params)
	}
	var res TradeHistoryResponse
	if err := c.doSigned(ctx, "tradeHistory", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) GetOpenOrders(ctx context.Context, pair string) (*OpenOrdersResponse, error) {
	params := map[string]string{}
	if pair != "" {
		params["pair"] = pair
	}
	var res OpenOrdersResponse
	if err := c.doSigned(ctx, "openOrders", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) GetOrderHistory(ctx context.Context, pair string, count, from int) (*OrderHistoryResponse, error) {
	params := map[string]string{"pair": pair}
	if count > 0 {
		params["count"] = strconv.Itoa(count)
	}
	if from > 0 {
		params["from"] = strconv.Itoa(from)
	}
	var res OrderHistoryResponse
	if err := c.doSigned(ctx, "orderHistory", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) GetOrder(ctx context.Context, pair string, orderID int64) (*GetOrderResponse, error) {
	params := map[string]string{"pair": pair, "order_id": strconv.FormatInt(orderID, 10)}
	var res GetOrderResponse
	if err := c.doSigned(ctx, "getOrder", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) GetOrderByClientOrderID(ctx context.Context, clientID string) (*GetOrderResponse, error) {
	params := map[string]string{"client_order_id": clientID}
	var res GetOrderResponse
	if err := c.doSigned(ctx, "getOrderByClientOrderId", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) CancelOrder(ctx context.Context, pair, tradeType string, orderID int64) (*CancelOrderResponse, error) {
	params := map[string]string{"pair": pair, "type": tradeType, "order_id": strconv.FormatInt(orderID, 10)}
	var res CancelOrderResponse
	if err := c.doSigned(ctx, "cancelOrder", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) CancelByClientOrderID(ctx context.Context, clientID string) (*CancelOrderResponse, error) {
	params := map[string]string{"client_order_id": clientID}
	var res CancelOrderResponse
	if err := c.doSigned(ctx, "cancelByClientOrderId", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) GetWithdrawFee(ctx context.Context, currency, network string) (*WithdrawFeeResponse, error) {
	params := map[string]string{"currency": currency}
	if network != "" {
		params["network"] = network
	}
	var res WithdrawFeeResponse
	if err := c.doSigned(ctx, "withdrawFee", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) WithdrawCoin(ctx context.Context, currency, network, address, amount, memo, reqID string) (*WithdrawCoinResponse, error) {
	params := map[string]string{"currency": currency, "network": network, "withdraw_address": address, "withdraw_amount": amount, "request_id": reqID}
	if memo != "" {
		params["withdraw_memo"] = memo
	}
	var res WithdrawCoinResponse
	if err := c.doSigned(ctx, "withdrawCoin", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) WithdrawByUsername(ctx context.Context, currency, amount, reqID, username string) (*WithdrawCoinResponse, error) {
	params := map[string]string{"currency": currency, "withdraw_input_method": "username", "withdraw_username": username, "withdraw_amount": amount, "request_id": reqID}
	var res WithdrawCoinResponse
	if err := c.doSigned(ctx, "withdrawCoin", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) ListDownline(ctx context.Context, page, limit int) (*ListDownlineResponse, error) {
	params := map[string]string{"page": strconv.Itoa(page), "limit": strconv.Itoa(limit)}
	var res ListDownlineResponse
	if err := c.doSigned(ctx, "listDownline", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) CheckDownline(ctx context.Context, email string) (*CheckDownlineResponse, error) {
	params := map[string]string{"email": email}
	var res CheckDownlineResponse
	if err := c.doSigned(ctx, "checkDownline", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *IndodaxClient) CreateVoucher(ctx context.Context, amount int, toEmail string) (*CreateVoucherResponse, error) {
	params := map[string]string{"amount": strconv.Itoa(amount), "to_email": toEmail}
	var res CreateVoucherResponse
	if err := c.doSigned(ctx, "createVoucher", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
