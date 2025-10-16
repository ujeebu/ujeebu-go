package ujeebu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// RawScrapeResponse represents the raw HTTP response from the scrape endpoint
type RawScrapeResponse struct {
	Body       []byte      // Raw response body
	StatusCode int         // HTTP status code
	Headers    http.Header // Response headers
}

// ContentType returns the Content-Type header value
func (r *RawScrapeResponse) ContentType() string {
	return r.Headers.Get("Content-Type")
}

// String returns the body as a string
func (r *RawScrapeResponse) String() string {
	return string(r.Body)
}

// ScrapeResponse represents the structured JSON response when `json=true`
type ScrapeResponse struct {
	Success         bool        `json:"success"`
	HTMLSource      string      `json:"html_source,omitempty"`
	HTML            string      `json:"html,omitempty"`
	PDF             string      `json:"pdf,omitempty"`
	Screenshot      string      `json:"screenshot,omitempty"`
	Result          any         `json:"result,omitempty"` // For extract_rules results
	StatusCode      int         `json:"status_code,omitempty"`
	ResponseHeaders http.Header `json:"response_headers,omitempty"`
}

// Scrape calls the Ujeebu Scrape API and returns structured JSON response
func (c *Client) Scrape(params ScrapeParams) (*ScrapeResponse, int, error) {
	// Force JSON output for structured response
	params.JSONOutput = true

	rawResp, credits, err := c.ScrapeWithContext(context.Background(), params)
	if err != nil {
		return nil, credits, err
	}

	// Parse JSON response into ScrapeResponse
	var scrapeResp ScrapeResponse
	if err := json.Unmarshal(rawResp.Body, &scrapeResp); err != nil {
		return nil, credits, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Populate HTTP metadata
	scrapeResp.StatusCode = rawResp.StatusCode
	scrapeResp.ResponseHeaders = rawResp.Headers

	return &scrapeResp, credits, nil
}

// ScrapeWithContext calls the Ujeebu Scrape API with context support and returns raw response
func (c *Client) ScrapeWithContext(ctx context.Context, params ScrapeParams) (*RawScrapeResponse, int, error) {
	// Validate required parameters
	if params.URL == "" {
		return nil, 0, &ValidationError{
			Field:   "URL",
			Message: "URL is required",
		}
	}

	// If FastMode is true, set Mode to d15de7
	if params.FastMode {
		params.Mode = "d15de7"
	}

	// Encode fields that need Base64
	if params.CustomJS != "" {
		params.CustomJS = encodeBase64(params.CustomJS)
	}
	if shouldEncodeWaitFor(params.WaitFor) {
		params.WaitFor = encodeBase64(params.WaitFor)
	}
	if params.ScrollCallback != "" {
		params.ScrollCallback = encodeBase64(params.ScrollCallback)
	}

	req := c.newRequest(ctx)

	// Add custom headers (prefixed with "UJB-")
	for key, value := range params.CustomHeaders {
		req.SetHeader("UJB-"+key, value)
	}

	req.SetError(&APIError{})
	var resp *resty.Response
	var err error

	if params.ExtractRules != nil {
		req.SetBody(params)
		req.SetHeader("Content-Type", "application/json")
		resp, err = req.Post("/scrape")
	} else {
		req.SetQueryParamsFromValues(params.toMap())
		resp, err = req.Get("/scrape")
	}

	if err != nil {
		return nil, 0, &NetworkError{Err: err}
	}

	if resp.IsError() {
		apiErr := resp.Error().(*APIError)
		apiErr.StatusCode = resp.StatusCode()
		return nil, 0, apiErr
	}

	rawResp := &RawScrapeResponse{
		Body:       resp.Body(),
		StatusCode: resp.StatusCode(),
		Headers:    resp.Header(),
	}

	return rawResp, getUjeebuCreditsFromResponse(resp), nil
}

// Screenshot retrieves the screenshot of the page with optional parameters.
func (c *Client) Screenshot(params ScrapeParams, fullPage bool, selector string) (string, int, error) {
	params.ResponseType = "screenshot"
	params.ScreenshotFullPage = fullPage
	params.ScreenshotPartial = selector

	response, credits, err := c.Scrape(params)
	if err != nil {
		return "", credits, err
	}
	return response.Screenshot, credits, nil
}

// PDF retrieves the PDF of the page with the specified scrape options.
func (c *Client) PDF(params ScrapeParams) (string, int, error) {
	params.ResponseType = "pdf"

	response, credits, err := c.Scrape(params)
	if err != nil {
		return "", credits, err
	}
	return response.PDF, credits, nil
}

// HTML retrieves the HTML content of the page with the specified scrape options.
func (c *Client) HTML(params ScrapeParams) (string, int, error) {
	params.ResponseType = "html"

	response, credits, err := c.Scrape(params)
	if err != nil {
		return "", credits, err
	}
	return response.HTML, credits, nil
}

// Raw retrieves the raw HTML source of the page with the specified scrape options.
func (c *Client) Raw(params ScrapeParams) (string, int, error) {
	params.ResponseType = "html"

	response, credits, err := c.Scrape(params)
	if err != nil {
		return "", credits, err
	}
	return response.HTMLSource, credits, nil
}
