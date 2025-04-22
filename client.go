package ujeebu

import (
	"errors"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client represents the Ujeebu API client
type Client struct {
	apiKey string
	client *resty.Client
}

var defaultBaseURL = "https://api.ujeebu.com"

// NewClient initializes and returns a new Ujeebu API client
func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("API key is required to initialize the Ujeebu client")
	}
	bu := os.Getenv("UJEEBU_BASE_URL")
	if bu == "" {
		bu = defaultBaseURL
	}

	return &Client{
		apiKey: apiKey,
		client: resty.New().
			SetBaseURL(bu).
			SetHeader("ApiKey", apiKey).
			SetHeader("User-Agent", "Ujeebu-GoSDK/1.0").
			SetTimeout(90 * time.Second),
	}, nil
}

// SetTimeout allows changing the client's timeout dynamically
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.SetTimeout(timeout)
}

// SetAPIKey allows setting or changing the API key dynamically
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
	c.client.SetHeader("ApiKey", apiKey)
}
