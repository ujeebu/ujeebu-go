package ujeebu

import (
	"context"
	"encoding/json"
	"fmt"
)

type ResponseMetadata struct {
	GoogleUrl       string `json:"google_url"`
	NumberOfResults int    `json:"number_of_results"`
	QueryDisplayed  string `json:"query_displayed"`
	ResultsTime     string `json:"results_time"`
}

type ResponsePagination struct {
	Google struct {
		Current    string `json:"current"`
		Next       string `json:"next"`
		OtherPages struct {
			Field1 string `json:"3"`
			Field2 string `json:"4"`
			Field3 string `json:"5"`
			Field4 string `json:"6"`
			Field5 string `json:"7"`
			Field6 string `json:"8"`
		} `json:"other_pages"`
	} `json:"google"`
	Api struct {
		Current    string `json:"current"`
		Next       string `json:"next"`
		OtherPages struct {
			Field1 string `json:"3"`
			Field2 string `json:"4"`
			Field3 string `json:"5"`
			Field4 string `json:"6"`
			Field5 string `json:"7"`
			Field6 string `json:"8"`
		} `json:"other_pages"`
	} `json:"api"`
}

type BaseResponse struct {
	Metadata   ResponseMetadata   `json:"metadata"`
	Pagination ResponsePagination `json:"pagination"`
}

// KnowledgeGraph represents the structure for the knowledge graph data
type KnowledgeGraph struct {
	Born      string `json:"born,omitempty"`
	Died      string `json:"died,omitempty"`
	Education string `json:"education,omitempty"`
	Height    string `json:"height,omitempty"`
	Parents   string `json:"parents,omitempty"`
	Siblings  string `json:"siblings,omitempty"`
	Title     string `json:"title,omitempty"`
	Type      string `json:"type,omitempty"`
}

// OrganicResult represents the structure for individual organic search results
type OrganicResult struct {
	Cite        string `json:"cite,omitempty"`
	Link        string `json:"link,omitempty"`
	Position    int    `json:"position,omitempty"`
	SiteName    string `json:"site_name,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// SearchVideo represents the structure for video search results
type SearchVideo struct {
	Author string `json:"author,omitempty"`
	Date   string `json:"date,omitempty"`
	Link   string `json:"link,omitempty"`
	Title  string `json:"title,omitempty"`
}

type TopStory struct {
	Link        string `json:"link,omitempty"`
	SiteName    string `json:"siteName,omitempty"`
	Description string `json:"description,omitempty"`
	PubDate     string `json:"pubDate,omitempty"`
}

// GoogleSearchResult extends BaseResponse to include additional JSON fields
type GoogleSearchResult struct {
	BaseResponse
	KnowledgeGraph   KnowledgeGraph  `json:"knowledge_graph"`
	OrganicResults   []OrganicResult `json:"organic_results"`
	RelatedQuestions []string        `json:"related_questions,omitempty"`
	TopStories       []TopStory      `json:"top_stories,omitempty"`
	Videos           []SearchVideo   `json:"videos,omitempty"`
}

type News struct {
	Date        string `json:"date,omitempty"`
	Description string `json:"description,omitempty"`
	Favicon     string `json:"favicon,omitempty"`
	Image       string `json:"image,omitempty"`
	Link        string `json:"link,omitempty"`
	Position    int    `json:"position,omitempty"`
	SiteName    string `json:"siteName,omitempty"`
	Title       string `json:"title,omitempty"`
}

// GoogleNewsResult represents metadata for a news search result with pagination
type GoogleNewsResult struct {
	BaseResponse
	News []News `json:"news,omitempty"`
}

type GoogleVideo struct {
	Author      string `json:"author"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Position    int    `json:"position"`
	Provider    string `json:"provider"`
	Site        string `json:"site"`
	Summary     string `json:"summary"`
	Title       string `json:"title"`
	Url         string `json:"url"`
}

type GoogleVideosResult struct {
	BaseResponse
	Videos []GoogleVideo `json:"videos,omitempty"`
}

type GoogleImage struct {
	GoogleThumbnail string `json:"google_thumbnail"`
	Height          int    `json:"height"`
	Image           string `json:"image"`
	Link            string `json:"link"`
	Position        int    `json:"position"`
	Source          string `json:"source"`
	Title           string `json:"title"`
	Width           int    `json:"width"`
}

type GoogleImageSuggestion struct {
	GoogleLink string `json:"google_link"`
	Position   int    `json:"position"`
	Thumbnail  string `json:"thumbnail"`
	Title      string `json:"title"`
	UjeebuLink string `json:"ujeebu_link"`
}

type GoogleImagesResult struct {
	BaseResponse
	Images      []GoogleImage           `json:"images,omitempty"`
	Suggestions []GoogleImageSuggestion `json:"suggestions,omitempty"`
}

type GoogleMap struct {
	Address      string      `json:"address"`
	Category     string      `json:"category"`
	Cid          string      `json:"cid"`
	OpeningHours interface{} `json:"opening_hours"`
	Position     int         `json:"position"`
	Rating       float64     `json:"rating"`
	Reviews      int         `json:"reviews"`
	Title        string      `json:"title"`
}

type GoogleMapsResult struct {
	BaseResponse
	Maps []GoogleMap `json:"maps_results,omitempty"`
}

// SerpParams represents the parameters used for the SERP API
type SerpParams struct {
	Search       string `json:"search,omitempty"`        // Search query
	URL          string `json:"url,omitempty"`           // URL of the search page
	SearchType   string `json:"search_type,omitempty"`   // "text", "images", "news", etc.
	Lang         string `json:"lang,omitempty"`          // Language of the results
	Location     string `json:"location,omitempty"`      // Location of search origin (e.g., "us", "uk")
	Device       string `json:"device,omitempty"`        // "desktop", "mobile", or "tablet"
	ResultsCount int    `json:"results_count,omitempty"` // Max results per page
	Page         int    `json:"page,omitempty"`          // Specific page to retrieve
	ExtraParams  string `json:"extra_params,omitempty"`  // Custom query parameters (&safe=active)
}

// Serp retrieves search results from Google using the SERP API and returns the processed data
func (c *Client) Serp(params SerpParams) ([]byte, int, error) {
	return c.SerpWithContext(context.Background(), params)
}

// SerpWithContext retrieves search results with context support
func (c *Client) SerpWithContext(ctx context.Context, params SerpParams) ([]byte, int, error) {
	// Validate required parameters - at least one of Search or URL must be provided
	if params.Search == "" && params.URL == "" {
		return nil, 0, &ValidationError{
			Field:   "Search/URL",
			Message: "Either Search or URL parameter is required",
		}
	}

	req := c.newRequest(ctx)
	req.SetQueryParams(serpParamsToMap(params))
	req.SetError(&APIError{})

	resp, err := req.Get("/serp")
	if err != nil {
		return nil, 0, &NetworkError{Err: err}
	}

	if resp.IsError() {
		apiErr := resp.Error().(*APIError)
		apiErr.StatusCode = resp.StatusCode()
		return nil, 0, apiErr
	}

	return resp.Body(), getUjeebuCreditsFromResponse(resp), nil
}

// Helper function to convert SerpParams into a query parameter map
func serpParamsToMap(params SerpParams) map[string]string {
	queryParams := map[string]string{}

	if params.Search != "" {
		queryParams["search"] = params.Search
	}
	if params.URL != "" {
		queryParams["url"] = params.URL
	}
	if params.SearchType != "" {
		queryParams["search_type"] = params.SearchType
	} else {
		queryParams["search_type"] = ""
	}
	if params.Lang != "" {
		queryParams["lang"] = params.Lang
	}
	if params.Location != "" {
		queryParams["location"] = params.Location
	}
	if params.Device != "" {
		queryParams["device"] = params.Device
	}
	if params.ResultsCount > 0 {
		queryParams["results_count"] = fmt.Sprintf("%d", params.ResultsCount)
	}
	if params.Page > 0 {
		queryParams["page"] = fmt.Sprintf("%d", params.Page)
	}
	if params.ExtraParams != "" {
		queryParams["extra_params"] = params.ExtraParams
	}

	return queryParams
}

// GoogleSearch Method for performing a Google Search
func (c *Client) GoogleSearch(params SerpParams) (GoogleSearchResult, int, error) {
	response, credits, err := c.Serp(params)
	if err != nil {
		return GoogleSearchResult{}, 0, err
	}

	// Deserialize JSON response "results" into specific type
	var results GoogleSearchResult
	err = json.Unmarshal(response, &results)
	if err != nil {
		return GoogleSearchResult{}, 0, fmt.Errorf("failed to parse Google Search results: %w", err)
	}

	return results, credits, nil
}

// GoogleImageSearch Method for performing a Google GoogleImage Search
func (c *Client) GoogleImageSearch(params SerpParams) (GoogleImagesResult, int, error) {
	params.SearchType = "images"
	response, credits, err := c.Serp(params)
	if err != nil {
		return GoogleImagesResult{}, 0, err
	}

	var results GoogleImagesResult
	err = json.Unmarshal(response, &results)
	if err != nil {
		return results, 0, fmt.Errorf("failed to parse Google GoogleImage results: %w", err)
	}

	return results, credits, nil
}

// GoogleNewsSearch Method for performing a Google News Search
func (c *Client) GoogleNewsSearch(params SerpParams) (GoogleNewsResult, int, error) {
	params.SearchType = "news"
	response, credits, err := c.Serp(params)
	if err != nil {
		return GoogleNewsResult{}, 0, err
	}

	var results GoogleNewsResult
	err = json.Unmarshal(response, &results)
	if err != nil {
		return results, 0, fmt.Errorf("failed to parse Google News results: %w", err)
	}

	return results, credits, nil
}

// GoogleVideoSearch Method for performing a Google SearchVideo Search
func (c *Client) GoogleVideoSearch(params SerpParams) (GoogleVideosResult, int, error) {
	params.SearchType = "videos"
	response, credits, err := c.Serp(params)
	if err != nil {
		return GoogleVideosResult{}, 0, err
	}

	var results GoogleVideosResult
	err = json.Unmarshal(response, &results)
	if err != nil {
		return results, 0, fmt.Errorf("failed to parse Google SearchVideo results: %w", err)
	}

	return results, credits, nil
}

// GoogleMapSearch Method for performing a Google Maps Search
func (c *Client) GoogleMapSearch(params SerpParams) (GoogleMapsResult, int, error) {
	params.SearchType = "maps"
	response, credits, err := c.Serp(params)
	if err != nil {
		return GoogleMapsResult{}, 0, err
	}

	var results GoogleMapsResult
	err = json.Unmarshal(response, &results)
	if err != nil {
		return results, 0, fmt.Errorf("failed to parse Google Map results: %w", err)
	}

	return results, credits, nil
}
