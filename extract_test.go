package ujeebu

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock server for Extract API
func setupMockExtractServer(response string, headers map[string]string, statusCode int) (*httptest.Server, *Client) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for k, v := range headers {
			w.Header().Add(k, v)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().SetBaseURL(mockServer.URL).SetTimeout(300 * time.Second),
	}

	return mockServer, client
}

func TestExtract_Success(t *testing.T) {
	mockResponse := `{
"article": {
		"url": "https://example.com/article",
		"title": "Sample Title",
		"text": "Sample extracted text",
		"author": "John Doe",
		"pub_date": "2023-01-01"
	}
}`

	mockServer, client := setupMockExtractServer(mockResponse, map[string]string{
		CreditsHeader: "10",
	}, http.StatusOK)
	defer mockServer.Close()

	params := ExtractParams{
		URL:    "https://example.com/article",
		Text:   true,
		Author: true,
	}

	article, credits, err := client.Extract(params)
	require.NoError(t, err)
	assert.NotNil(t, article)
	assert.Equal(t, "Sample Title", article.Title)
	assert.Equal(t, "John Doe", article.Author)
	assert.Equal(t, "Sample extracted text", article.Text)
	assert.Equal(t, "2023-01-01", article.PubDate)
	assert.Equal(t, 10, credits)
}

func TestExtract_ErrorResponse(t *testing.T) {
	mockResponse := `{
		"url": "https://example.com/article",
		"message": "Invalid API Key",
		"error_code": "401"
	}`

	mockServer, client := setupMockExtractServer(mockResponse, map[string]string{}, http.StatusUnauthorized)
	defer mockServer.Close()

	params := ExtractParams{URL: "https://example.com/article"}
	article, credits, err := client.Extract(params)

	assert.Error(t, err)
	assert.Nil(t, article)
	assert.Equal(t, 0, credits)
	assert.Contains(t, err.Error(), "Invalid API Key")
}

func TestExtract_Timeout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"url": "https://example.com", "title": "Delayed Response"}`))
	}))
	defer mockServer.Close()

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().SetBaseURL(mockServer.URL).SetTimeout(1 * time.Second),
	}

	params := ExtractParams{URL: "https://example.com"}
	article, credits, err := client.Extract(params)

	assert.Error(t, err)
	assert.Nil(t, article)
	assert.Equal(t, 0, credits)
}

func TestExtract_Validation_URLRequired_EvenWithRawHTML(t *testing.T) {
	client := &Client{}

	_, _, err := client.Extract(ExtractParams{RawHTML: "<html></html>"})
	require.Error(t, err)

	var validationErr *ValidationError
	require.ErrorAs(t, err, &validationErr)
	assert.Equal(t, "URL", validationErr.Field)
}

func TestExtract_RawHTML_UsesPOSTAndBody(t *testing.T) {
	const rawHTML = "<html><body>ok</body></html>"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/extract", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(b, &payload))
		assert.Equal(t, "https://example.com/article", payload["url"])
		assert.Equal(t, rawHTML, payload["raw_html"])

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set(CreditsHeader, "7")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"article":{"title":"ok"}}`))
	}))
	defer mockServer.Close()

	client, err := NewClient("test_api_key", WithBaseURL(mockServer.URL), WithRetry(0, 0, 0))
	require.NoError(t, err)

	article, credits, err := client.Extract(ExtractParams{
		URL:     "https://example.com/article",
		RawHTML: rawHTML,
	})
	require.NoError(t, err)
	require.NotNil(t, article)
	assert.Equal(t, "ok", article.Title)
	assert.Equal(t, 7, credits)
}

func TestArticle_Unmarshal_ToleratesBooleans(t *testing.T) {
	var r ExtractResponse
	data := []byte(`{"article":{"title":"t","summary":false,"author":false,"text":null}}`)
	require.NoError(t, json.Unmarshal(data, &r))
	require.NotNil(t, r.Article)
	assert.Equal(t, "t", r.Article.Title)
	assert.Equal(t, "", r.Article.Summary)
	assert.Equal(t, "", r.Article.Author)
	assert.Equal(t, "", r.Article.Text)
}
