# Indodax API Client

A comprehensive Go client for the Indodax cryptocurrency exchange API. This client provides full access to all private API endpoints as documented in the official Indodax API documentation.

## Features

- ✅ Complete HMAC-SHA512 signature authentication
- ✅ All private API endpoints implemented
- ✅ Proper error handling and rate limiting detection
- ✅ Thread-safe nonce management
- ✅ Comprehensive response structures
- ✅ Helper methods for common trading operations
- ✅ Support for all order types (limit, market, stop-limit)
- ✅ Withdrawal and deposit management
- ✅ Transaction history tracking

## Installation

```bash
go get github.com/frederickmarvel/supernova/internal/client
```

## Quick Start

### 1. Set up your API credentials

First, you need to create an API key from your Indodax account:

1. Log in to your Indodax account
2. Go to API Management
3. Create a new API key with appropriate permissions (view, trade, withdraw)
4. Save your API key and secret

### 2. Initialize the client

```go
package main

import (
    "log"
    "os"
    "github.com/frederickmarvel/supernova/internal/client"
)

func main() {
    apiKey := os.Getenv("INDODAX_API_KEY")
    apiSecret := os.Getenv("INDODAX_API_SECRET")
    
    if apiKey == "" || apiSecret == "" {
        log.Fatal("Please set INDODAX_API_KEY and INDODAX_API_SECRET environment variables")
    }
    
    indodax := client.NewIndodaxClient(apiKey, apiSecret)
    
    // Use the client...
}
```

### 3. Basic usage examples

#### Get account information
```go
info, err := indodax.GetInfo()
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("User ID: %s\n", info.Return.UserID)
fmt.Printf("Name: %s\n", info.Return.Name)
fmt.Printf("Email: %s\n", info.Return.Email)
```

#### Get balance for a specific currency
```go
btcBalance, err := indodax.GetBalance("btc")
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("BTC Balance: %f\n", btcBalance)
```

#### Place a limit buy order
```go
response, err := indodax.PlaceLimitBuy("btc_idr", 50000000, 0.001, "my-order-123", "GTC")
if err != nil {
    log.Printf("Error placing order: %v", err)
    return
}

fmt.Printf("Order placed! Order ID: %v\n", response.Return["order_id"])
```

## API Endpoints

### Account Information

#### GetInfo()
Retrieves user account information including balances, addresses, and account details.

```go
info, err := indodax.GetInfo()
```

### Trading

#### Trade()
Executes a trade order with full parameter control.

```go
response, err := indodax.Trade("btc_idr", "limit", "buy", 50000000, 0.001, 0, "order-123", "GTC")
```

#### Helper Methods for Trading

```go
// Limit orders
response, err := indodax.PlaceLimitBuy("btc_idr", 50000000, 0.001, "order-123", "GTC")
response, err := indodax.PlaceLimitSell("btc_idr", 51000000, 0.001, "order-124", "GTC")

// Market orders
response, err := indodax.PlaceMarketBuy("btc_idr", 1000000, "order-125")  // Buy with IDR
response, err := indodax.PlaceMarketSell("btc_idr", 0.001, "order-126")   // Sell BTC
```

### Order Management

#### GetOpenOrders()
Retrieves current open orders.

```go
orders, err := indodax.GetOpenOrders("btc_idr")
```

#### GetOrderHistory()
Retrieves order history.

```go
history, err := indodax.GetOrderHistory("btc_idr", 100, 0)
```

#### GetOrder()
Retrieves a specific order by order ID.

```go
order, err := indodax.GetOrder("btc_idr", 12345678)
```

#### GetOrderByClientOrderID()
Retrieves a specific order by client order ID.

```go
order, err := indodax.GetOrderByClientOrderID("my-order-123")
```

#### CancelOrder()
Cancels an existing order.

```go
response, err := indodax.CancelOrder("btc_idr", "limit", "buy", 12345678)
```

#### CancelOrderByClientOrderID()
Cancels an existing order by client order ID.

```go
response, err := indodax.CancelOrderByClientOrderID("my-order-123")
```

### Trade History

#### GetTradeHistory()
Retrieves trade history with various filtering options.

```go
history, err := indodax.GetTradeHistory("btc_idr", 100, 0, 0, "desc", "", "", 0)
```

### Transaction History

#### GetTransactionHistory()
Retrieves deposit and withdrawal history.

```go
history, err := indodax.GetTransactionHistory("2024-01-01", "2024-01-07")
```

### Withdrawal Operations

#### GetWithdrawFee()
Gets withdrawal fee for a currency.

```go
fee, err := indodax.GetWithdrawFee("btc", "")
```

#### WithdrawCoin()
Withdraws cryptocurrency to an external address.

```go
response, err := indodax.WithdrawCoin("btc", "mainnet", "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "0.001", "", "req-123")
```

#### WithdrawCoinByUsername()
Withdraws cryptocurrency to an Indodax username.

```go
response, err := indodax.WithdrawCoinByUsername("btc", "0.001", "", "req-124", "username123")
```

### Downline Management

#### ListDownline()
Lists all downline users.

```go
downline, err := indodax.ListDownline(1, 10)
```

#### CheckDownline()
Checks if an email exists in current user's downline.

```go
result, err := indodax.CheckDownline("user@example.com")
```

### Voucher Creation (Partner Only)

#### CreateVoucher()
Creates an Indodax voucher (requires partner status).

```go
voucher, err := indodax.CreateVoucher(10000, "recipient@example.com")
```

## Error Handling

The client provides comprehensive error handling:

```go
response, err := indodax.GetInfo()
if err != nil {
    // Handle error
    log.Printf("API Error: %v", err)
    return
}

// Check for API-specific errors
if response.Success == 0 {
    log.Printf("API returned error: %s (code: %s)", response.Error, response.ErrorCode)
    return
}
```

## Rate Limiting

The client automatically detects rate limiting (HTTP 429) and returns appropriate errors:

```go
response, err := indodax.Trade("btc_idr", "limit", "buy", 50000000, 0.001, 0, "order-123", "GTC")
if err != nil {
    if strings.Contains(err.Error(), "rate limit exceeded") {
        // Handle rate limiting
        log.Println("Rate limit hit, wait before retrying")
        time.Sleep(5 * time.Second)
    }
    return
}
```

## Security Features

- **HMAC-SHA512 Signing**: All requests are properly signed using HMAC-SHA512
- **Nonce Management**: Thread-safe nonce generation to prevent replay attacks
- **Request Validation**: Proper parameter validation and encoding
- **Error Handling**: Comprehensive error handling for API responses

## Configuration

### Environment Variables

Set these environment variables for your API credentials:

```bash
export INDODAX_API_KEY="your-api-key"
export INDODAX_API_SECRET="your-api-secret"
```

### API Permissions

Make sure your API key has the appropriate permissions:

- **view**: Required for reading account info, balances, orders, and history
- **trade**: Required for placing and canceling orders
- **withdraw**: Required for withdrawal operations

## Running the Example

```bash
cd supernova
go run examples/indodax_example.go
```

Make sure to set your API credentials as environment variables first.

## API Documentation Reference

This client implements all endpoints from the official Indodax Private API documentation:

- Base URL: `https://indodax.com/tapi`
- Authentication: HMAC-SHA512 signature
- Request Method: POST
- Content Type: `application/x-www-form-urlencoded`

For detailed API documentation, refer to the official Indodax API docs.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License. 