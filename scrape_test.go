package ujeebu

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock server for Scrape API
func setupMockScrapeServer(response string, headers map[string]string, contentType string, statusCode int) (*httptest.Server, *Client) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range headers {
			w.Header().Add(k, v)
		}
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().SetBaseURL(mockServer.URL).SetTimeout(10 * time.Second),
	}

	return mockServer, client
}

func TestScrape_Success(t *testing.T) {
	mockResponse := `{
		"success": true,
		"html": "<html><body>Sample HTML</body></html>"
	}`

	mockServer, client := setupMockScrapeServer(mockResponse, map[string]string{
		"ujb-credits": "5",
	}, "application/json", http.StatusOK)
	defer mockServer.Close()

	params := ScrapeParams{
		URL:          "https://example.com",
		ResponseType: "html",
		JS:           true,
	}

	scraped, credits, err := client.Scrape(params)
	require.NoError(t, err)
	assert.NotNil(t, scraped)
	assert.Equal(t, 5, credits)
	assert.Contains(t, scraped.HTML, "Sample HTML")
}

func TestScrape_ErrorResponse(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com",
		"message": "Page Not Found",
		"error_code": "404"
	}`

	mockServer, client := setupMockScrapeServer(mockResponse, map[string]string{}, "application/json", http.StatusNotFound)
	defer mockServer.Close()

	params := ScrapeParams{URL: "https://example.com"}
	scraped, credits, err := client.Scrape(params)

	assert.Error(t, err)
	assert.Nil(t, scraped)
	assert.Equal(t, 0, credits)
	assert.Contains(t, err.Error(), "Page Not Found")
}

func TestScrape_Timeout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "html": "<p>Delayed Response</p>"}`))
	}))
	defer mockServer.Close()

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().SetBaseURL(mockServer.URL).SetTimeout(1 * time.Second),
	}

	params := ScrapeParams{URL: "https://example.com"}
	scraped, c, err := client.Scrape(params)

	assert.Error(t, err)
	assert.Equal(t, 0, c)
	assert.Nil(t, scraped)
}

func TestScreenshot(t *testing.T) {
	tests := []struct {
		name            string
		mockResponse    string
		mockHeaders     map[string]string
		contentType     string
		statusCode      int
		params          ScrapeParams
		fullPage        bool
		selector        string
		expectedError   string
		expectedImage   string
		expectedCredits int
	}{
		{
			name: "Success, full-page screenshot",
			mockResponse: `{
				"success": true,
				"screenshot": "base64-screenshot-data"
			}`,
			mockHeaders: map[string]string{
				"ujb-credits": "10",
			},
			contentType:     "application/json",
			statusCode:      http.StatusOK,
			params:          ScrapeParams{URL: "https://example.com"},
			fullPage:        true,
			selector:        "",
			expectedError:   "",
			expectedImage:   "base64-screenshot-data",
			expectedCredits: 10,
		},
		{
			name: "Partial screenshot with selector",
			mockResponse: `{
				"success": true,
				"screenshot": "base64-partial-screenshot-data"
			}`,
			mockHeaders: map[string]string{
				"ujb-credits": "15",
			},
			contentType:     "application/json",
			statusCode:      http.StatusOK,
			params:          ScrapeParams{URL: "https://example.com"},
			fullPage:        false,
			selector:        "#target-element",
			expectedError:   "",
			expectedImage:   "base64-partial-screenshot-data",
			expectedCredits: 15,
		},
		{
			name: "Error response",
			mockResponse: `{
				"url": "https://example.com",
				"message": "Invalid API Key",
				"error_code": "401"
			}`,
			mockHeaders:     nil,
			contentType:     "application/json",
			statusCode:      http.StatusUnauthorized,
			params:          ScrapeParams{URL: "https://example.com"},
			fullPage:        true,
			selector:        "",
			expectedError:   "Invalid API Key",
			expectedImage:   "",
			expectedCredits: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(time.Duration(0))
				for k, v := range tt.mockHeaders {
					w.Header().Add(k, v)
				}
				if tt.contentType != "" {
					w.Header().Set("Content-Type", tt.contentType)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			client := &Client{
				apiKey: "test_api_key",
				client: resty.New().SetBaseURL(mockServer.URL).SetTimeout(1 * time.Second),
			}

			image, credits, err := client.Screenshot(tt.params, tt.fullPage, tt.selector)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, "", image)
				assert.Equal(t, tt.expectedCredits, credits)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedImage, image)
				assert.Equal(t, tt.expectedCredits, credits)
			}
		})
	}
}
