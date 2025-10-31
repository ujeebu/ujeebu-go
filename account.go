package ujeebu

import (
	"context"
	"fmt"
)

// AccountResponse represents the response from the Ujeebu Account API
type AccountResponse struct {
	Balance             int     `json:"balance"`
	DaysTillNextBilling int     `json:"days_till_next_billing"`
	NextBillingDate     *string `json:"next_billing_date"` // Pointer to handle null values
	Plan                string  `json:"plan"`
	Quota               int     `json:"quota"`
	ConcurrentRequests  int     `json:"concurrent_requests"`
	TotalRequests       int     `json:"total_requests"`
	RequestPerSecond    int     `json:"requests_per_second"`
	Used                int     `json:"used"`
	UsedPercent         float64 `json:"used_percent"`
	UserID              string  `json:"userid"`
}

// Account retrieves account information including usage and billing details
func (c *Client) Account() (*AccountResponse, error) {
	return c.AccountWithContext(context.Background())
}

// AccountWithContext retrieves account information with context support
func (c *Client) AccountWithContext(ctx context.Context) (*AccountResponse, error) {
	req := c.newRequest(ctx)
	req.SetResult(&AccountResponse{}).SetError(&APIError{})

	resp, err := req.Get("/account")
	if err != nil {
		return nil, &NetworkError{Err: err}
	}

	if resp.IsError() {
		apiErr := resp.Error().(*APIError)
		apiErr.StatusCode = resp.StatusCode()
		return nil, apiErr
	}

	res := resp.Result()
	if r, ok := res.(*AccountResponse); ok {
		return r, nil
	}
	return nil, fmt.Errorf("account API response is not a valid AccountResponse")
}
