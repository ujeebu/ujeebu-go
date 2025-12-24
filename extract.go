package ujeebu

import (
	"bytes"
	"context"
	"encoding/json"
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

// UnmarshalJSON makes Article tolerant to inconsistent API payloads.
// Some endpoints may return boolean/null for fields that are usually strings.
func (a *Article) UnmarshalJSON(data []byte) error {
	type articleWire struct {
		URL          json.RawMessage `json:"url"`
		CanonicalURL json.RawMessage `json:"canonical_url"`
		Title        json.RawMessage `json:"title"`
		Text         json.RawMessage `json:"text"`
		HTML         json.RawMessage `json:"html"`
		Summary      json.RawMessage `json:"summary"`
		Image        json.RawMessage `json:"image"`
		Images       []string        `json:"images"`
		Media        []string        `json:"media"`
		Language     json.RawMessage `json:"language"`
		Author       json.RawMessage `json:"author"`
		PubDate      json.RawMessage `json:"pub_date"`
		ModifiedDate json.RawMessage `json:"modified_date"`
		SiteName     json.RawMessage `json:"site_name"`
		Favicon      json.RawMessage `json:"favicon"`
		Encoding     json.RawMessage `json:"encoding"`
		IsArticle    json.RawMessage `json:"is_article,omitempty"`
		Pages        []string        `json:"pages"`
	}

	var w articleWire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}

	a.URL = rawJSONToString(w.URL)
	a.CanonicalURL = rawJSONToString(w.CanonicalURL)
	a.Title = rawJSONToString(w.Title)
	a.Text = rawJSONToString(w.Text)
	a.HTML = rawJSONToString(w.HTML)
	a.Summary = rawJSONToString(w.Summary)
	a.Image = rawJSONToString(w.Image)
	a.Images = w.Images
	a.Media = w.Media
	a.Language = rawJSONToString(w.Language)
	a.Author = rawJSONToString(w.Author)
	a.PubDate = rawJSONToString(w.PubDate)
	a.ModifiedDate = rawJSONToString(w.ModifiedDate)
	a.SiteName = rawJSONToString(w.SiteName)
	a.Favicon = rawJSONToString(w.Favicon)
	a.Encoding = rawJSONToString(w.Encoding)
	a.IsArticle = rawJSONToFloat64(w.IsArticle)
	a.Pages = w.Pages

	return nil
}

func rawJSONToString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}

	// Common fast-paths
	if string(raw) == "null" || string(raw) == "false" || string(raw) == "true" {
		return ""
	}

	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}

	// Fallback: try to decode as a number and stringify it
	var n json.Number
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&n); err == nil {
		return n.String()
	}

	return ""
}

func rawJSONToFloat64(raw json.RawMessage) float64 {
	if len(raw) == 0 || string(raw) == "null" {
		return 0
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return f
	}
	return 0
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
