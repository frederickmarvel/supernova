package client

import (
	"strconv"
)

// Base response structure for all API calls
type BaseResponse struct {
	Success int    `json:"success"`
	Error   string `json:"error,omitempty"`
}

// GetInfoResponse represents the response from getInfo endpoint
type GetInfoResponse struct {
	BaseResponse
	Return struct {
		ServerTime         int64                  `json:"server_time"`
		Balance            map[string]string      `json:"balance"`
		BalanceHold        map[string]string      `json:"balance_hold"`
		Address            map[string]string      `json:"address"`
		Network            map[string]interface{} `json:"network"`
		MemoIsRequired     map[string]interface{} `json:"memo_is_required"`
		UserID             string                 `json:"user_id"`
		Name               string                 `json:"name"`
		Email              string                 `json:"email"`
		ProfilePicture     interface{}            `json:"profile_picture"`
		VerificationStatus string                 `json:"verification_status"`
		GauthEnable        bool                   `json:"gauth_enable"`
		WithdrawStatus     int                    `json:"withdraw_status"`
	} `json:"return"`
}

// TransHistoryResponse represents the response from transHistory endpoint
type TransHistoryResponse struct {
	BaseResponse
	Return struct {
		Withdraw map[string][]Transaction `json:"withdraw"`
		Deposit  map[string][]Transaction `json:"deposit"`
	} `json:"return"`
}

// Transaction represents a single transaction
type Transaction struct {
	Status      string `json:"status"`
	Type        string `json:"type"`
	RP          string `json:"rp,omitempty"`
	Fee         string `json:"fee,omitempty"`
	Amount      string `json:"amount,omitempty"`
	SubmitTime  string `json:"submit_time"`
	SuccessTime string `json:"success_time"`
	ID          string `json:"withdraw_id,omitempty"`
	TX          string `json:"tx"`
	BTC         string `json:"btc,omitempty"`
}

// TradeResponse represents the response from trade endpoint
type TradeResponse struct {
	BaseResponse
	Return struct {
		OrderID       string `json:"order_id"`
		ClientOrderID string `json:"client_order_id,omitempty"`
		ReceiveIDR    string `json:"receive_idr,omitempty"`
		ReceiveBTC    string `json:"receive_btc,omitempty"`
		SpentIDR      string `json:"spent_idr,omitempty"`
		SpentBTC      string `json:"spent_btc,omitempty"`
		Fee           string `json:"fee,omitempty"`
		SubmitTime    string `json:"submit_time"`
		FinishTime    string `json:"finish_time,omitempty"`
		Status        string `json:"status"`
	} `json:"return"`
}

// TradeHistoryOption is a function type for trade history options
type TradeHistoryOption func(map[string]string)

// WithCount sets the count parameter for trade history
func WithCount(count int) TradeHistoryOption {
	return func(params map[string]string) {
		params["count"] = strconv.Itoa(count)
	}
}

// WithFrom sets the from parameter for trade history
func WithFrom(from int) TradeHistoryOption {
	return func(params map[string]string) {
		params["from"] = strconv.Itoa(from)
	}
}

// WithOrder sets the order parameter for trade history
func WithOrder(order string) TradeHistoryOption {
	return func(params map[string]string) {
		params["order"] = order
	}
}

// TradeHistoryResponse represents the response from tradeHistory endpoint
type TradeHistoryResponse struct {
	BaseResponse
	Return struct {
		Trades []TradeHistoryItem `json:"trades"`
	} `json:"return"`
}

// TradeHistoryItem represents a single trade history item
type TradeHistoryItem struct {
	OrderID       string `json:"order_id"`
	ClientOrderID string `json:"client_order_id,omitempty"`
	Type          string `json:"type"`
	Price         string `json:"price"`
	SubmitTime    string `json:"submit_time"`
	FinishTime    string `json:"finish_time"`
	Status        string `json:"status"`
	OrderBTC      string `json:"order_btc,omitempty"`
	RemainBTC     string `json:"remain_btc,omitempty"`
	OrderIDR      string `json:"order_idr,omitempty"`
	RemainIDR     string `json:"remain_idr,omitempty"`
	ReceiveIDR    string `json:"receive_idr,omitempty"`
	ReceiveBTC    string `json:"receive_btc,omitempty"`
	SpentIDR      string `json:"spent_idr,omitempty"`
	SpentBTC      string `json:"spent_btc,omitempty"`
	Fee           string `json:"fee,omitempty"`
}

// OpenOrdersResponse represents the response from openOrders endpoint
type OpenOrdersResponse struct {
	BaseResponse
	Return struct {
		Orders []OpenOrder `json:"orders"`
	} `json:"return"`
}

// OpenOrder represents a single open order
type OpenOrder struct {
	OrderID       string `json:"order_id"`
	ClientOrderID string `json:"client_order_id,omitempty"`
	SubmitTime    string `json:"submit_time"`
	Price         string `json:"price"`
	Type          string `json:"type"`
	OrderType     string `json:"order_type"`
	OrderBTC      string `json:"order_btc,omitempty"`
	RemainBTC     string `json:"remain_btc,omitempty"`
	OrderIDR      string `json:"order_idr,omitempty"`
	RemainIDR     string `json:"remain_idr,omitempty"`
}

// OrderHistoryResponse represents the response from orderHistory endpoint
type OrderHistoryResponse struct {
	BaseResponse
	Return struct {
		Orders []OrderHistoryItem `json:"orders"`
	} `json:"return"`
}

// OrderHistoryItem represents a single order history item
type OrderHistoryItem struct {
	OrderID       string `json:"order_id"`
	ClientOrderID string `json:"client_order_id,omitempty"`
	Type          string `json:"type"`
	Price         string `json:"price"`
	SubmitTime    string `json:"submit_time"`
	FinishTime    string `json:"finish_time"`
	Status        string `json:"status"`
	OrderBTC      string `json:"order_btc,omitempty"`
	RemainBTC     string `json:"remain_btc,omitempty"`
	OrderIDR      string `json:"order_idr,omitempty"`
	RemainIDR     string `json:"remain_idr,omitempty"`
	ReceiveIDR    string `json:"receive_idr,omitempty"`
	ReceiveBTC    string `json:"receive_btc,omitempty"`
	SpentIDR      string `json:"spent_idr,omitempty"`
	SpentBTC      string `json:"spent_btc,omitempty"`
	Fee           string `json:"fee,omitempty"`
}

// GetOrderResponse represents the response from getOrder endpoint
type GetOrderResponse struct {
	BaseResponse
	Return struct {
		Order OrderDetail `json:"order"`
	} `json:"return"`
}

// OrderDetail represents a single order detail
type OrderDetail struct {
	OrderID       string `json:"order_id"`
	Price         string `json:"price"`
	Type          string `json:"type"`
	OrderRP       string `json:"order_rp,omitempty"`
	RemainRP      string `json:"remain_rp,omitempty"`
	SubmitTime    string `json:"submit_time"`
	FinishTime    string `json:"finish_time"`
	Status        string `json:"status"`
	ReceiveIDR    string `json:"receive_idr,omitempty"`
	ClientOrderID string `json:"client_order_id,omitempty"`
}

// CancelOrderResponse represents the response from cancelOrder endpoint
type CancelOrderResponse struct {
	BaseResponse
	Return struct {
		OrderID       string `json:"order_id"`
		ClientOrderID string `json:"client_order_id,omitempty"`
		Type          string `json:"type"`
		Price         string `json:"price"`
		SubmitTime    string `json:"submit_time"`
		FinishTime    string `json:"finish_time"`
		Status        string `json:"status"`
		OrderBTC      string `json:"order_btc,omitempty"`
		RemainBTC     string `json:"remain_btc,omitempty"`
		OrderIDR      string `json:"order_idr,omitempty"`
		RemainIDR     string `json:"remain_idr,omitempty"`
		ReceiveIDR    string `json:"receive_idr,omitempty"`
		ReceiveBTC    string `json:"receive_btc,omitempty"`
		SpentIDR      string `json:"spent_idr,omitempty"`
		SpentBTC      string `json:"spent_btc,omitempty"`
		Fee           string `json:"fee,omitempty"`
	} `json:"return"`
}

// WithdrawFeeResponse represents the response from withdrawFee endpoint
type WithdrawFeeResponse struct {
	BaseResponse
	Return struct {
		Fee string `json:"fee"`
	} `json:"return"`
}

// WithdrawCoinResponse represents the response from withdrawCoin endpoint
type WithdrawCoinResponse struct {
	BaseResponse
	Status           string `json:"status"`
	WithdrawCurrency string `json:"withdraw_currency"`
	WithdrawAddress  string `json:"withdraw_address"`
	WithdrawAmount   string `json:"withdraw_amount"`
	Fee              string `json:"fee"`
	AmountAfterFee   string `json:"amount_after_fee"`
	SubmitTime       string `json:"submit_time"`
	WithdrawID       string `json:"withdraw_id"`
	TXID             string `json:"txid"`
	WithdrawUsername string `json:"withdraw_username,omitempty"`
}

// ListDownlineResponse represents the response from listDownline endpoint
type ListDownlineResponse struct {
	BaseResponse
	Return struct {
		Downlines []Downline `json:"downlines"`
		Total     int        `json:"total"`
	} `json:"return"`
}

// Downline represents a single downline
type Downline struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	JoinDate string `json:"join_date"`
}

// CheckDownlineResponse represents the response from checkDownline endpoint
type CheckDownlineResponse struct {
	BaseResponse
	Return struct {
		IsDownline bool `json:"is_downline"`
	} `json:"return"`
}

// CreateVoucherResponse represents the response from createVoucher endpoint
type CreateVoucherResponse struct {
	BaseResponse
	Return struct {
		VoucherID string `json:"voucher_id"`
		Amount    int    `json:"amount"`
		ToEmail   string `json:"to_email"`
	} `json:"return"`
}
