package ujeebu

import (
	"fmt"
)

type AccountResponse struct {
	Balance             float64 `json:"balance"`
	DaysTillNextBilling int     `json:"days_till_next_billing"`
	NextBillingDate     string  `json:"next_billing_date"`
	Plan                string  `json:"plan"`
	Quota               string  `json:"quota"`
	ConcurrentRequests  int     `json:"concurrent_requests"`
	TotalRequests       string  `json:"total_requests"`
	Used                string  `json:"used"`
	UsedPercent         float64 `json:"used_percent"`
	UserID              string  `json:"userid"`
}

func (c *Client) Account() (*AccountResponse, error) {
	req := c.client.R()

	req = req.SetResult(&AccountResponse{})

	resp, err := req.Get("/account")

	fmt.Println(string(resp.Body()))

	if err != nil {
		return nil, fmt.Errorf("account API error: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("account API error: %v", resp.Error())
	}

	res := resp.Result()
	if r, ok := res.(*AccountResponse); ok {
		return r, nil
	}
	return nil, fmt.Errorf("account API response is not a valid AccountResponse")
}
