package ujeebu

import (
	"context"
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
	IsArticle    float64  `json:"is_article,omitempty"`
	Pages        []string `json:"pages"`
}

// ExtractResponse represents the full response from the Extract API
type ExtractResponse struct {
	Article    *Article `json:"article"`
	Time       float64  `json:"time,omitempty"`
	JS         bool     `json:"js,omitempty"`
	Pagination bool     `json:"pagination,omitempty"`
}

// Extract calls the Ujeebu Extract API and returns structured data
func (c *Client) Extract(params ExtractParams) (*Article, int, error) {
	return c.ExtractWithContext(context.Background(), params)
}

// ExtractWithContext calls the Ujeebu Extract API with context support
func (c *Client) ExtractWithContext(ctx context.Context, params ExtractParams) (*Article, int, error) {
	// Validate required parameters
	if params.URL == "" && params.RawHTML == "" {
		return nil, 0, &ValidationError{
			Field:   "URL",
			Message: "URL or RawHTML is required",
		}
	}

	// If FastMode is true, set Mode to d15de7
	if params.FastMode {
		params.Mode = "d15de7"
	}

	req := c.newRequest(ctx)

	// Add custom headers (prefixed with "UJB-")
	for key, value := range params.CustomHeaders {
		req.SetHeader("UJB-"+key, value)
	}

	req.SetResult(&ExtractResponse{}).SetError(&APIError{})
	var resp *resty.Response
	var err error

	if params.RawHTML != "" {
		req.SetBody(params)
		req.SetHeader("Content-Type", "application/json")
		resp, err = req.Post("/extract")
	} else {
		req.SetQueryParamsFromValues(params.toMap())
		resp, err = req.Get("/extract")
	}

	if err != nil {
		return nil, 0, &NetworkError{Err: err}
	}

	if resp.IsError() {
		apiErr := resp.Error().(*APIError)
		apiErr.StatusCode = resp.StatusCode()
		return nil, 0, apiErr
	}

	res := resp.Result()
	if r, ok := res.(*ExtractResponse); ok && r.Article != nil {
		return r.Article, getUjeebuCreditsFromResponse(resp), nil
	}
	return nil, 0, fmt.Errorf("extract API response is not a valid ExtractResponse")
}
