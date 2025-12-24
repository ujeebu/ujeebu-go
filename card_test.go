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

// Mock server for Card API
func setupMockCardServer(response string, headers map[string]string, statusCode int) (*httptest.Server, *Client) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("ApiKey") != "test_api_key" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message": "Invalid API Key"}`))
			return
		}

		// Add custom headers
		for k, v := range headers {
			w.Header().Add(k, v)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().
			SetBaseURL(mockServer.URL).
			SetHeader("ApiKey", "test_api_key").
			SetTimeout(5 * time.Second),
	}

	return mockServer, client
}

func TestCard_Success(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com/article",
		"lang": "en",
		"favicon": "https://example.com/favicon.ico",
		"title": "Example Article Title",
		"summary": "This is a brief summary of the article content.",
		"author": "John Doe",
		"date_published": "2024-01-15T10:30:00Z",
		"date_modified": "2024-01-16T14:20:00Z",
		"image": "https://example.com/image.jpg",
		"site_name": "Example Site",
		"charset": "UTF-8",
		"keywords": ["technology", "web", "development"],
		"time": 0.234
	}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{
		CreditsHeader: "5",
	}, http.StatusOK)
	defer mockServer.Close()

	params := CardParams{
		URL: "https://example.com/article",
	}

	card, credits, err := client.Card(params)
	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "https://example.com/article", card.URL)
	assert.Equal(t, "en", card.Lang)
	assert.Equal(t, "https://example.com/favicon.ico", card.Favicon)
	assert.Equal(t, "Example Article Title", card.Title)
	assert.Equal(t, "This is a brief summary of the article content.", card.Summary)
	assert.Equal(t, "John Doe", card.Author)
	assert.Equal(t, "2024-01-15T10:30:00Z", card.DatePublished)
	assert.Equal(t, "2024-01-16T14:20:00Z", card.DateModified)
	assert.Equal(t, "https://example.com/image.jpg", card.Image)
	assert.Equal(t, "Example Site", card.SiteName)
	assert.Equal(t, "UTF-8", card.Charset)
	assert.Equal(t, []string{"technology", "web", "development"}, card.Keywords)
	assert.Equal(t, 0.234, card.Time)
	assert.Equal(t, 5, credits)
}

func TestCard_SuccessWithJS(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com/spa",
		"title": "SPA Title",
		"summary": "Single page application content",
		"time": 1.456
	}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{
		CreditsHeader: "8",
	}, http.StatusOK)
	defer mockServer.Close()

	params := CardParams{
		URL:       "https://example.com/spa",
		JS:        true,
		JSTimeout: 5000,
	}

	card, credits, err := client.Card(params)
	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "https://example.com/spa", card.URL)
	assert.Equal(t, "SPA Title", card.Title)
	assert.Equal(t, "Single page application content", card.Summary)
	assert.Equal(t, 1.456, card.Time)
	assert.Equal(t, 8, credits)
}

func TestCard_WithCustomHeaders(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify custom headers are prefixed with UJB-
		assert.Equal(t, "custom-value", r.Header.Get("UJB-Custom-Header"))
		assert.Equal(t, "another-value", r.Header.Get("UJB-X-Custom"))

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set(CreditsHeader, "3")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"url": "https://example.com", "title": "Test"}`))
	}))
	defer mockServer.Close()

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().
			SetBaseURL(mockServer.URL).
			SetHeader("ApiKey", "test_api_key").
			SetTimeout(5 * time.Second),
	}

	params := CardParams{
		URL: "https://example.com",
		CustomHeaders: map[string]string{
			"Custom-Header": "custom-value",
			"X-Custom":      "another-value",
		},
	}

	card, credits, err := client.Card(params)
	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "Test", card.Title)
	assert.Equal(t, 3, credits)
}

func TestCard_WithProxy(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com",
		"title": "Proxied Request",
		"summary": "Content fetched through proxy"
	}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{
		CreditsHeader: "6",
	}, http.StatusOK)
	defer mockServer.Close()

	params := CardParams{
		URL:          "https://example.com",
		ProxyType:    "residential",
		ProxyCountry: "us",
	}

	card, credits, err := client.Card(params)
	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "Proxied Request", card.Title)
	assert.Equal(t, 6, credits)
}

func TestCard_WithSession(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com",
		"title": "Session Request",
		"summary": "Content with session"
	}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{
		CreditsHeader: "4",
	}, http.StatusOK)
	defer mockServer.Close()

	params := CardParams{
		URL:       "https://example.com",
		SessionID: "test-session-123",
	}

	card, credits, err := client.Card(params)
	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "Session Request", card.Title)
	assert.Equal(t, 4, credits)
}

func TestCard_MissingURL(t *testing.T) {
	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().
			SetBaseURL("https://api.example.com").
			SetHeader("ApiKey", "test_api_key").
			SetTimeout(5 * time.Second),
	}

	params := CardParams{
		// URL is intentionally missing
	}

	card, credits, err := client.Card(params)
	require.Error(t, err)
	assert.Nil(t, card)
	assert.Equal(t, 0, credits)

	// Check if it's a ValidationError
	var validationErr *ValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, "URL", validationErr.Field)
	assert.Contains(t, validationErr.Message, "required")
}

func TestCard_ErrorResponse(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com/error",
		"message": "Failed to fetch URL",
		"error_code": "FETCH_ERROR"
	}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{}, http.StatusBadRequest)
	defer mockServer.Close()

	params := CardParams{URL: "https://example.com/error"}
	card, credits, err := client.Card(params)

	require.Error(t, err)
	assert.Nil(t, card)
	assert.Equal(t, 0, credits)

	// Check if it's an APIError
	var apiErr *APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "Failed to fetch URL")
}

func TestCard_UnauthorizedError(t *testing.T) {
	mockResponse := `{
		"message": "Invalid API Key",
		"error_code": "UNAUTHORIZED"
	}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{}, http.StatusUnauthorized)
	defer mockServer.Close()

	params := CardParams{URL: "https://example.com"}
	card, credits, err := client.Card(params)

	require.Error(t, err)
	assert.Nil(t, card)
	assert.Equal(t, 0, credits)

	// Check if it's an APIError with Unauthorized status
	var apiErr *APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.True(t, apiErr.IsUnauthorized())
	assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
}

func TestCard_NotFoundError(t *testing.T) {
	mockResponse := `{
		"message": "URL not found",
		"error_code": "NOT_FOUND"
	}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{}, http.StatusNotFound)
	defer mockServer.Close()

	params := CardParams{URL: "https://example.com/nonexistent"}
	card, credits, err := client.Card(params)

	require.Error(t, err)
	assert.Nil(t, card)
	assert.Equal(t, 0, credits)

	// Check if it's an APIError with NotFound status
	var apiErr *APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.True(t, apiErr.IsNotFound())
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
}

func TestCard_RateLimitError(t *testing.T) {
	mockResponse := `{
		"message": "Rate limit exceeded",
		"error_code": "RATE_LIMIT"
	}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{}, http.StatusTooManyRequests)
	defer mockServer.Close()

	params := CardParams{URL: "https://example.com"}
	card, credits, err := client.Card(params)

	require.Error(t, err)
	assert.Nil(t, card)
	assert.Equal(t, 0, credits)

	// Check if it's an APIError with RateLimited status
	var apiErr *APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.True(t, apiErr.IsRateLimited())
	assert.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
}

func TestCard_Timeout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"url": "https://example.com", "title": "Delayed Response"}`))
	}))
	defer mockServer.Close()

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().
			SetBaseURL(mockServer.URL).
			SetHeader("ApiKey", "test_api_key").
			SetTimeout(500 * time.Millisecond),
	}

	params := CardParams{URL: "https://example.com"}
	card, credits, err := client.Card(params)

	require.Error(t, err)
	assert.Nil(t, card)
	assert.Equal(t, 0, credits)

	// Check if it's a NetworkError
	var netErr *NetworkError
	assert.ErrorAs(t, err, &netErr)
}

func TestCard_InvalidJSON(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Return invalid JSON
		_, _ = w.Write([]byte(`{invalid json`))
	}))
	defer mockServer.Close()

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().
			SetBaseURL(mockServer.URL).
			SetHeader("ApiKey", "test_api_key").
			SetTimeout(5 * time.Second),
	}

	params := CardParams{URL: "https://example.com"}
	card, credits, err := client.Card(params)

	// Should handle invalid JSON gracefully
	require.Error(t, err)
	assert.Nil(t, card)
	assert.Equal(t, 0, credits)
}

func TestCard_EmptyResponse(t *testing.T) {
	mockResponse := `{}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{
		CreditsHeader: "2",
	}, http.StatusOK)
	defer mockServer.Close()

	params := CardParams{URL: "https://example.com"}
	card, credits, err := client.Card(params)

	require.NoError(t, err)
	assert.NotNil(t, card)
	// All fields should be empty/zero values
	assert.Equal(t, "", card.URL)
	assert.Equal(t, "", card.Title)
	assert.Equal(t, "", card.Summary)
	assert.Equal(t, 2, credits)
}

func TestCardWithContext_Success(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com",
		"title": "Context Test",
		"summary": "Testing with context"
	}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{
		CreditsHeader: "3",
	}, http.StatusOK)
	defer mockServer.Close()

	ctx := context.Background()
	params := CardParams{URL: "https://example.com"}

	card, credits, err := client.CardWithContext(ctx, params)
	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "Context Test", card.Title)
	assert.Equal(t, 3, credits)
}

func TestCardWithContext_Cancellation(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"url": "https://example.com", "title": "Slow Response"}`))
	}))
	defer mockServer.Close()

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().
			SetBaseURL(mockServer.URL).
			SetHeader("ApiKey", "test_api_key").
			SetTimeout(10 * time.Second),
	}

	// Create a context that will be cancelled
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	params := CardParams{URL: "https://example.com"}
	card, credits, err := client.CardWithContext(ctx, params)

	require.Error(t, err)
	assert.Nil(t, card)
	assert.Equal(t, 0, credits)

	// Check if it's a NetworkError (context cancellation wraps as network error)
	var netErr *NetworkError
	assert.ErrorAs(t, err, &netErr)
}

func TestCard_NilContext(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com",
		"title": "Nil Context Test"
	}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{
		CreditsHeader: "1",
	}, http.StatusOK)
	defer mockServer.Close()

	params := CardParams{URL: "https://example.com"}

	// Pass nil context - should default to Background
	card, credits, err := client.CardWithContext(context.TODO(), params)
	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "Nil Context Test", card.Title)
	assert.Equal(t, 1, credits)
}

func TestCard_AllParameters(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com",
		"title": "All Parameters Test",
		"summary": "Testing all parameters"
	}`

	mockServer, client := setupMockCardServer(mockResponse, map[string]string{
		CreditsHeader: "10",
	}, http.StatusOK)
	defer mockServer.Close()

	params := CardParams{
		URL:          "https://example.com",
		JS:           true,
		Timeout:      30,
		JSTimeout:    10000,
		ProxyType:    "datacenter",
		ProxyCountry: "de",
		CustomProxy:  "http://proxy.example.com:8080",
		AutoProxy:    true,
		SessionID:    "session-abc-123",
		CustomHeaders: map[string]string{
			"X-Custom-Header": "test-value",
		},
	}

	card, credits, err := client.Card(params)
	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "All Parameters Test", card.Title)
	assert.Equal(t, 10, credits)
}

func TestCard_NoCreditsHeader(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com",
		"title": "No Credits Test"
	}`

	// Don't set the credits header
	mockServer, client := setupMockCardServer(mockResponse, map[string]string{}, http.StatusOK)
	defer mockServer.Close()

	params := CardParams{URL: "https://example.com"}
	card, credits, err := client.Card(params)

	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "No Credits Test", card.Title)
	assert.Equal(t, 0, credits) // Should default to 0 when header is missing
}
