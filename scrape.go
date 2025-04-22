package ujeebu

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// ScrapeResponse represents the structured JSON response when `json=true`
type ScrapeResponse struct {
	Success    bool   `json:"success"`
	HTMLSource string `json:"html_source,omitempty"`
	HTML       string `json:"html,omitempty"`
	PDF        string `json:"pdf,omitempty"`
	Screenshot string `json:"screenshot,omitempty"`
}

// ScrapeError represents an error response from the API
type ScrapeError struct {
	URL       string   `json:"url"`
	Message   string   `json:"message"`
	ErrorCode int      `json:"error_code"`
	Errors    []string `json:"errors"`
}

// Scrape calls the Ujeebu Scrape API and returns the requested web page data
func (c *Client) Scrape(params ScrapeParams) (response *ScrapeResponse, credits int, err error) {

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
	params.JSONOutput = true

	req := c.client.R()

	// Add custom headers (prefixed with "UJB-")
	for key, value := range params.CustomHeaders {
		req.SetHeader("UJB-"+key, value)
	}

	req.SetResult(&ScrapeResponse{}).SetError(&ScrapeError{})
	var resp *resty.Response
	if params.ExtractRules != nil {
		req.SetBody(params)
		req.SetHeader("Content-Type", "application/json")
		resp, err = req.Post("/scrape")
	} else {
		req.SetQueryParamsFromValues(params.toMap())
		resp, err = req.Get("/scrape")
	}

	if err != nil {
		return nil, 0, fmt.Errorf("scrape API request failed: %w", err)
	}

	if resp.IsError() {
		apiErr := resp.Error().(*ScrapeError)
		return nil, 0, fmt.Errorf("scrape API error: %s (%d)", apiErr.Message, apiErr.ErrorCode)
	}

	return resp.Result().(*ScrapeResponse), getUjeebuCreditsFromResponse(resp), nil
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
