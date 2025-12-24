package ujeebu

import (
	"context"
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
		"ujb-credits":     "5",
		"X-Custom-Header": "test-value",
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
	assert.Equal(t, http.StatusOK, scraped.StatusCode)
	assert.NotNil(t, scraped.ResponseHeaders)
	assert.Equal(t, "5", scraped.ResponseHeaders.Get("Ujb-Credits"))
	assert.Equal(t, "test-value", scraped.ResponseHeaders.Get("X-Custom-Header"))
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
		_, _ = w.Write([]byte(`{"success": true, "html": "<p>Delayed Response</p>"}`))
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

func TestScrapeWithContext_RawHTML(t *testing.T) {
	mockHTMLResponse := `<!DOCTYPE html>
<html>
<head><title>Test Page</title></head>
<body><h1>Hello World</h1></body>
</html>`

	mockServer, client := setupMockScrapeServer(mockHTMLResponse, map[string]string{
		"ujb-credits": "3",
		"X-Custom":    "value",
	}, "text/html", http.StatusOK)
	defer mockServer.Close()

	params := ScrapeParams{
		URL: "https://example.com",
		JS:  false,
	}

	rawResp, credits, err := client.ScrapeWithContext(context.Background(), params)
	require.NoError(t, err)
	assert.NotNil(t, rawResp)
	assert.Equal(t, 3, credits)
	assert.Equal(t, http.StatusOK, rawResp.StatusCode)
	assert.Equal(t, "text/html", rawResp.ContentType())
	assert.Contains(t, rawResp.String(), "Hello World")
	assert.Equal(t, mockHTMLResponse, string(rawResp.Body))
	assert.Equal(t, "value", rawResp.Headers.Get("X-Custom"))
}

func TestScrapeWithContext_RawImage(t *testing.T) {
	// Simulate binary image data
	mockImageData := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46} // JPEG header

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("ujb-credits", "5")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockImageData)
	}))
	defer mockServer.Close()

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().SetBaseURL(mockServer.URL).SetTimeout(10 * time.Second),
	}

	params := ScrapeParams{
		URL: "https://example.com/image.jpg",
	}

	rawResp, credits, err := client.ScrapeWithContext(context.Background(), params)
	require.NoError(t, err)
	assert.NotNil(t, rawResp)
	assert.Equal(t, 5, credits)
	assert.Equal(t, http.StatusOK, rawResp.StatusCode)
	assert.Equal(t, "image/jpeg", rawResp.ContentType())
	assert.Equal(t, mockImageData, rawResp.Body)
}

func TestScrapeWithContext_JSON(t *testing.T) {
	mockJSONResponse := `{"success": true, "data": "test"}`

	mockServer, client := setupMockScrapeServer(mockJSONResponse, map[string]string{
		"ujb-credits": "2",
	}, "application/json", http.StatusOK)
	defer mockServer.Close()

	params := ScrapeParams{
		URL:        "https://api.example.com",
		JSONOutput: true,
	}

	rawResp, credits, err := client.ScrapeWithContext(context.Background(), params)
	require.NoError(t, err)
	assert.NotNil(t, rawResp)
	assert.Equal(t, 2, credits)
	assert.Equal(t, "application/json", rawResp.ContentType())
	assert.JSONEq(t, mockJSONResponse, rawResp.String())
}

func TestScrapeWithContext_ErrorResponse(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com",
		"message": "Invalid URL",
		"error_code": "400"
	}`

	mockServer, client := setupMockScrapeServer(mockResponse, map[string]string{}, "application/json", http.StatusBadRequest)
	defer mockServer.Close()

	params := ScrapeParams{URL: "https://example.com"}
	rawResp, credits, err := client.ScrapeWithContext(context.Background(), params)

	assert.Error(t, err)
	assert.Nil(t, rawResp)
	assert.Equal(t, 0, credits)
	assert.Contains(t, err.Error(), "Invalid URL")
}

func TestScrapeWithContext_ValidationError(t *testing.T) {
	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().SetBaseURL("http://test.com").SetTimeout(10 * time.Second),
	}

	params := ScrapeParams{} // Missing URL

	rawResp, credits, err := client.ScrapeWithContext(context.Background(), params)

	assert.Error(t, err)
	assert.Nil(t, rawResp)
	assert.Equal(t, 0, credits)
	assert.IsType(t, &ValidationError{}, err)
	assert.Contains(t, err.Error(), "URL is required")
}

func TestScrapeWithContext_WithContext(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.Header().Set("ujb-credits", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html>Test</html>"))
	}))
	defer mockServer.Close()

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().SetBaseURL(mockServer.URL).SetTimeout(2 * time.Second),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	params := ScrapeParams{URL: "https://example.com"}
	rawResp, credits, err := client.ScrapeWithContext(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, rawResp)
	assert.Equal(t, 0, credits)
}

func TestScrape_ForcesJSONOutput(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that json parameter is set to true
		assert.Equal(t, "true", r.URL.Query().Get("json"))

		w.Header().Set("ujb-credits", "4")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true, "html": "<div>Test</div>"}`))
	}))
	defer mockServer.Close()

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().SetBaseURL(mockServer.URL).SetTimeout(10 * time.Second),
	}

	params := ScrapeParams{
		URL: "https://example.com",
		// JSONOutput not set - should be forced to true by Scrape()
	}

	scrapeResp, credits, err := client.Scrape(params)
	require.NoError(t, err)
	assert.NotNil(t, scrapeResp)
	assert.Equal(t, 4, credits)
	assert.True(t, scrapeResp.Success)
	assert.Contains(t, scrapeResp.HTML, "Test")
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
