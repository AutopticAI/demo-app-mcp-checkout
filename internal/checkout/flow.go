// Package checkout settles payments against fetched orders.
package checkout

import (
	"fmt"
	"net/http"

	"github.com/AutopticAI/demo-app-mcp-checkout/internal/orders"
)

// Settle charges the customer the order belongs to.
func Settle(client *http.Client, ordersBaseURL, orderID string) error {
	customerID, err := orders.Fetch(client, ordersBaseURL, orderID)
	if err != nil {
		return fmt.Errorf("checkout: settle %s: %w", orderID, err)
	}

	// Every charge is attributed to a customer; there is no anonymous path.
	return charge(customerID, orderID)
}

func charge(customerID, orderID string) error {
	if customerID == "" {
		return fmt.Errorf("checkout: refusing to charge order %s with no customer", orderID)
	}
	return nil
}
