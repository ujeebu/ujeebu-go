package ujeebu

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// DefaultBaseURL is the default Ujeebu API base URL
	DefaultBaseURL = "https://api.ujeebu.com"
	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 90 * time.Second
	// DefaultUserAgent is the default User-Agent header value
	DefaultUserAgent = "Ujeebu-GoSDK/2.0"
)

// Client represents the Ujeebu API client with support for all Ujeebu endpoints
type Client struct {
	apiKey    string
	baseURL   string
	client    *resty.Client
	debug     bool
	logger    Logger
	retryConf *RetryConfig
}

// Logger is an interface for logging
type Logger interface {
	Printf(format string, v ...interface{})
}

// RetryConfig defines retry behavior for failed requests
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts
	MaxRetries int
	// WaitTime is the initial wait time between retries
	WaitTime time.Duration
	// MaxWaitTime is the maximum wait time between retries
	MaxWaitTime time.Duration
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
		c.client.SetBaseURL(url)
	}
}

// WithTimeout sets a custom timeout for API requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.client.SetTimeout(timeout)
	}
}

// WithDebug enables debug mode with detailed logging
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
		c.client.SetDebug(debug)
	}
}

// WithLogger sets a custom logger for the client
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithRetry configures retry behavior for failed requests
func WithRetry(maxRetries int, waitTime, maxWaitTime time.Duration) ClientOption {
	return func(c *Client) {
		c.retryConf = &RetryConfig{
			MaxRetries:  maxRetries,
			WaitTime:    waitTime,
			MaxWaitTime: maxWaitTime,
		}
		c.client.
			SetRetryCount(maxRetries).
			SetRetryWaitTime(waitTime).
			SetRetryMaxWaitTime(maxWaitTime)
	}
}

// WithUserAgent sets a custom User-Agent header
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.client.SetHeader("User-Agent", userAgent)
	}
}

// NewClient creates a new Ujeebu API client with the provided API key and options
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("API key is required to initialize the Ujeebu client")
	}

	// Get base URL from environment or use default
	baseURL := os.Getenv("UJEEBU_BASE_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	// Create client with defaults
	client := &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: resty.New().
			SetBaseURL(baseURL).
			SetHeader("ApiKey", apiKey).
			SetHeader("User-Agent", DefaultUserAgent).
			SetTimeout(DefaultTimeout),
		logger: log.Default(),
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// SetTimeout allows changing the client's timeout dynamically
// Deprecated: Use WithTimeout option in NewClient instead
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.SetTimeout(timeout)
}

// SetAPIKey allows setting or changing the API key dynamically
// Deprecated: Create a new client with the new API key instead
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
	c.client.SetHeader("ApiKey", apiKey)
}

// GetAPIKey returns the current API key
func (c *Client) GetAPIKey() string {
	return c.apiKey
}

// GetBaseURL returns the current base URL
func (c *Client) GetBaseURL() string {
	return c.baseURL
}

// newRequest creates a new request with context support
func (c *Client) newRequest(ctx context.Context) *resty.Request {
	if ctx == nil {
		ctx = context.Background()
	}
	return c.client.R().SetContext(ctx)
}
