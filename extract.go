package ujeebu

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Article represents the extracted data returned by the API
type Article struct {
	URL          string   `json:"url"`
	CanonicalURL string   `json:"canonical_url"`
	Title        string   `json:"title"`
	Text         string   `json:"text"`
	HTML         string   `json:"html"`
	Summary      string   `json:"summary"`
	Image        string   `json:"image"`
	Images       []string `json:"images"`
	Media        []string `json:"media"`
	Language     string   `json:"language"`
	Author       string   `json:"author"`
	PubDate      string   `json:"pub_date"`
	ModifiedDate string   `json:"modified_date"`
	SiteName     string   `json:"site_name"`
	Favicon      string   `json:"favicon"`
	Encoding     string   `json:"encoding"`
	Pages        []string `json:"pages"`
	Time         float64  `json:"time"`
	JS           bool     `json:"js"`
	Pagination   bool     `json:"pagination"`
}
type ExtractResponse struct {
	Article *Article `json:"article"`
}

// ExtractError represents an error response from the API
type ExtractError struct {
	URL       string   `json:"url"`
	Message   string   `json:"message"`
	ErrorCode int      `json:"error_code"`
	Errors    []string `json:"errors"`
}

// Extract calls the Ujeebu Extract API and returns structured data
func (c *Client) Extract(params ExtractParams) (article *Article, credits int, err error) {
	req := c.client.R()

	// Add custom headers (prefixed with "UJB-")
	for key, value := range params.CustomHeaders {
		req.SetHeader("UJB-"+key, value)
	}

	req = req.SetResult(&ExtractResponse{}).SetError(&ExtractError{})
	var resp *resty.Response

	if params.RawHTML != "" {
		req.SetBody(params)
		req.SetHeader("Content-Type", "application/json")
		resp, err = req.Post("/extract")
	} else {
		req.SetQueryParamsFromValues(params.toMap())
		resp, err = req.Get("/extract")
	}

	if err != nil {
		return nil, 0, fmt.Errorf("extract API request failed: %w", err)
	}

	if resp.IsError() {
		apiErr := resp.Error().(*ExtractError)
		return nil, 0, fmt.Errorf("extract API error: %s (%d)", apiErr.Message, apiErr.ErrorCode)
	}

	res := resp.Result()
	if r, ok := res.(*ExtractResponse); ok {
		return r.Article, getUjeebuCreditsFromResponse(resp), nil
	}
	return nil, 0, fmt.Errorf("extract API response is not a valid ExtractResponse")
}
