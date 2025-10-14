package ujeebu

import (
	"context"
	"net/url"
)

// CardParams defines the parameters for the Card API (Article Preview API)
type CardParams struct {
	URL           string            `json:"url"`
	JS            bool              `json:"js,omitempty"`
	Timeout       int               `json:"timeout,omitempty"`
	JSTimeout     int               `json:"js_timeout,omitempty"`
	ProxyType     string            `json:"proxy_type,omitempty"`
	ProxyCountry  string            `json:"proxy_country,omitempty"`
	CustomProxy   string            `json:"custom_proxy,omitempty"`
	AutoProxy     bool              `json:"auto_proxy,omitempty"`
	SessionID     string            `json:"session_id,omitempty"`
	CustomHeaders map[string]string `json:"-"` // UJB-prefixed headers
}

// CardResponse represents the response from the Ujeebu Card API
type CardResponse struct {
	URL           string   `json:"url,omitempty"`
	Lang          string   `json:"lang,omitempty"`
	Favicon       string   `json:"favicon,omitempty"`
	Title         string   `json:"title,omitempty"`
	Summary       string   `json:"summary,omitempty"`
	Author        string   `json:"author,omitempty"`
	DatePublished string   `json:"date_published,omitempty"`
	DateModified  string   `json:"date_modified,omitempty"`
	Image         string   `json:"image,omitempty"`
	SiteName      string   `json:"site_name,omitempty"`
	Charset       string   `json:"charset,omitempty"`
	Keywords      []string `json:"keywords,omitempty"`
	Time          float64  `json:"time,omitempty"`
}

// Converts struct fields to query parameters
func (p CardParams) toMap() url.Values {
	return structToQueryParams(p)
}

// Card retrieves article card/preview information from a URL
// This is a fast alternative to Extract that mostly relies on meta tags
func (c *Client) Card(params CardParams) (*CardResponse, int, error) {
	return c.CardWithContext(context.Background(), params)
}

// CardWithContext retrieves article card/preview information with context support
func (c *Client) CardWithContext(ctx context.Context, params CardParams) (*CardResponse, int, error) {
	// Validate required parameters
	if params.URL == "" {
		return nil, 0, &ValidationError{
			Field:   "URL",
			Message: "URL is required",
		}
	}

	req := c.newRequest(ctx)

	// Add custom headers (prefixed with "UJB-")
	for key, value := range params.CustomHeaders {
		req.SetHeader("UJB-"+key, value)
	}

	// Set up response and error structs
	req.SetResult(&CardResponse{}).SetError(&APIError{})

	// Set query parameters
	req.SetQueryParamsFromValues(params.toMap())

	// Execute GET request
	resp, err := req.Get("/card")
	if err != nil {
		return nil, 0, &NetworkError{Err: err}
	}

	// Handle error responses
	if resp.IsError() {
		apiErr := resp.Error().(*APIError)
		apiErr.StatusCode = resp.StatusCode()
		return nil, 0, apiErr
	}

	// Extract credits from response header
	credits := getUjeebuCreditsFromResponse(resp)

	// Return successful response
	result := resp.Result().(*CardResponse)
	return result, credits, nil
}
