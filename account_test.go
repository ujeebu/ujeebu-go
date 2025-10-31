package ujeebu

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountInfo(t *testing.T) {
	nextBillingDate := "2025-12-01"
	tests := []struct {
		name         string
		mockResponse string
		mockStatus   int
		expectedErr  string
		expectedResp *AccountResponse
	}{
		{
			name:         "successful response",
			mockResponse: `{"balance":100,"days_till_next_billing":15,"next_billing_date":"2025-12-01","plan":"STARTER","quota":5000,"concurrent_requests":10,"total_requests":14,"requests_per_second":2,"used":95,"used_percent":1.9,"userid":"8155"}`,
			mockStatus:   http.StatusOK,
			expectedErr:  "",
			expectedResp: &AccountResponse{
				Balance:             100,
				DaysTillNextBilling: 15,
				NextBillingDate:     &nextBillingDate,
				Plan:                "STARTER",
				Quota:               5000,
				ConcurrentRequests:  10,
				TotalRequests:       14,
				RequestPerSecond:    2,
				Used:                95,
				UsedPercent:         1.9,
				UserID:              "8155",
			},
		},
		{
			name:         "successful response with null next_billing_date",
			mockResponse: `{"balance":50,"days_till_next_billing":0,"next_billing_date":null,"plan":"FREE","quota":1000,"concurrent_requests":5,"total_requests":100,"requests_per_second":0,"used":250,"used_percent":25.0,"userid":"1234"}`,
			mockStatus:   http.StatusOK,
			expectedErr:  "",
			expectedResp: &AccountResponse{
				Balance:             50,
				DaysTillNextBilling: 0,
				NextBillingDate:     nil,
				Plan:                "FREE",
				Quota:               1000,
				ConcurrentRequests:  5,
				TotalRequests:       100,
				RequestPerSecond:    0,
				Used:                250,
				UsedPercent:         25.0,
				UserID:              "1234",
			},
		},
		{
			name:         "API request failure",
			mockResponse: `{"message": "Not found"}`,
			mockStatus:   http.StatusNotFound,
			expectedErr:  "Not found",
			expectedResp: nil,
		},
		{
			name:         "error response from API",
			mockResponse: `{"message": "Internal server error"}`,
			mockStatus:   http.StatusInternalServerError,
			expectedErr:  "Internal server error",
			expectedResp: nil,
		},
		{
			name:         "invalid response type",
			mockResponse: "invalid-json",
			mockStatus:   http.StatusOK,
			expectedErr:  "invalid character",
			expectedResp: nil,
		},
		{
			name:         "unauthorized error",
			mockResponse: `{"message": "Invalid API Key"}`,
			mockStatus:   http.StatusUnauthorized,
			expectedErr:  "Invalid API Key",
			expectedResp: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.mockStatus)
				time.Sleep(time.Duration(0))
				_, _ = w.Write([]byte(tc.mockResponse))
			}))
			defer mockServer.Close()

			client := &Client{
				apiKey: "test_api_key",
				client: resty.New().SetBaseURL(mockServer.URL).SetTimeout(1 * time.Second),
			}
			resp, err := client.Account()

			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedResp, resp)
		})
	}
}
