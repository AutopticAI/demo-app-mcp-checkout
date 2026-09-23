// Package orders reads the Orders API contract.
package orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrMissingCustomer reports an order that arrived with no customer on it.
// Checkout cannot attribute a payment without one, so this is fatal to the
// flow rather than something to fall back from.
var ErrMissingCustomer = errors.New("orders: response carried no customer_id")

// orderPayload mirrors OrderResponse in demo-app-mcp-orders. Only the fields
// checkout actually reads are listed.
type orderPayload struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	TotalCents int    `json:"total_cents"`
	Currency   string `json:"currency"`
}

// Fetch reads one order and returns the customer it belongs to.
func Fetch(client *http.Client, baseURL, orderID string) (string, error) {
	resp, err := client.Get(fmt.Sprintf("%s/orders/%s", baseURL, orderID))
	if err != nil {
		return "", fmt.Errorf("orders: get %s: %w", orderID, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("orders: get %s: status %d", orderID, resp.StatusCode)
	}

	var payload orderPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("orders: decode %s: %w", orderID, err)
	}

	if payload.CustomerID == "" {
		return "", ErrMissingCustomer
	}
	return payload.CustomerID, nil
}
