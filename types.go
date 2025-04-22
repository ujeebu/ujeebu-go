package ujeebu

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// ExtractParams defines the parameters for the Extract API
type ExtractParams struct {
	URL                          string            `json:"url"`
	RawHTML                      string            `json:"raw_html,omitempty"`
	JS                           bool              `json:"js,omitempty"`
	Text                         bool              `json:"text,omitempty"`
	HTML                         bool              `json:"html,omitempty"`
	Media                        bool              `json:"media,omitempty"`
	Feeds                        bool              `json:"feeds,omitempty"`
	Images                       bool              `json:"images,omitempty"`
	Author                       bool              `json:"author,omitempty"`
	PubDate                      bool              `json:"pub_date,omitempty"`
	Partial                      string            `json:"partial,omitempty"`
	IsArticle                    bool              `json:"is_article,omitempty"`
	QuickMode                    bool              `json:"quick_mode,omitempty"`
	StripTags                    string            `json:"strip_tags,omitempty"`
	Timeout                      int               `json:"timeout,omitempty"`
	JSTimeout                    int               `json:"js_timeout,omitempty"`
	ScrollDown                   bool              `json:"scroll_down,omitempty"`
	ImageAnalysis                bool              `json:"image_analysis,omitempty"`
	MinImageWidth                int               `json:"min_image_width,omitempty"`
	MinImageHeight               int               `json:"min_image_height,omitempty"`
	ImageTimeout                 int               `json:"image_timeout,omitempty"`
	ReturnOnlyEnclosedTextImages bool              `json:"return_only_enclosed_text_images,omitempty"`
	ProxyType                    string            `json:"proxy_type,omitempty"`
	ProxyCountry                 string            `json:"proxy_country,omitempty"`
	CustomProxy                  string            `json:"custom_proxy,omitempty"`
	AutoProxy                    bool              `json:"auto_proxy,omitempty"`
	SessionID                    string            `json:"session_id,omitempty"`
	Pagination                   bool              `json:"pagination,omitempty"`
	PaginationMaxPages           int               `json:"pagination_max_pages,omitempty"`
	CustomHeaders                map[string]string `json:"-"` // UJB-prefixed headers
}

// Converts struct fields to query parameters
func (p ExtractParams) toMap() url.Values {
	return structToQueryParams(p)
}

// ScrapeParams defines the parameters for the Scrape API
type ScrapeParams struct {
	URL                 string            `json:"url"`
	ResponseType        string            `json:"response_type,omitempty"`
	JSONOutput          bool              `json:"json,omitempty"`
	UserAgent           string            `json:"useragent,omitempty"`
	Cookies             string            `json:"cookies,omitempty"`
	Timeout             int               `json:"timeout,omitempty"`
	JS                  bool              `json:"js"`
	JSTimeout           int               `json:"js_timeout,omitempty"`
	CustomJS            string            `json:"custom_js,omitempty"`
	WaitFor             string            `json:"wait_for,omitempty"`
	WaitForTimeout      int               `json:"wait_for_timeout,omitempty"`
	ScreenshotFullPage  bool              `json:"screenshot_fullpage,omitempty"`
	ScreenshotPartial   string            `json:"screenshot_partial,omitempty"`
	ScrollDown          bool              `json:"scroll_down,omitempty"`
	ScrollWait          int               `json:"scroll_wait,omitempty"`
	ProgressiveScroll   bool              `json:"progressive_scroll,omitempty"`
	ProxyType           string            `json:"proxy_type,omitempty"`
	ProxyCountry        string            `json:"proxy_country,omitempty"`
	CustomProxy         string            `json:"custom_proxy,omitempty"`
	CustomProxyUsername string            `json:"custom_proxy_username,omitempty"`
	CustomProxyPassword string            `json:"custom_proxy_password,omitempty"`
	AutoProxy           bool              `json:"auto_proxy,omitempty"`
	SessionID           string            `json:"session_id,omitempty"`
	ScrollCallback      string            `json:"scroll_callback,omitempty"`
	ScrollToSelector    string            `json:"scroll_to_selector,omitempty"`
	Device              string            `json:"device,omitempty"`
	WindowWidth         int               `json:"window_width,omitempty"`
	WindowHeight        int               `json:"window_height,omitempty"`
	BlockAds            bool              `json:"block_ads,omitempty"`
	BlockResources      bool              `json:"block_resources,omitempty"`
	ExtractRules        map[string]any    `json:"extract_rules,omitempty"`
	StripTags           string            `json:"strip_tags,omitempty"`
	HTTPMethod          string            `json:"http_method,omitempty"`
	PostData            string            `json:"post_data,omitempty"`
	Mode                string            `json:"mode,omitempty"`
	CustomHeaders       map[string]string `json:"-"`
}

// Converts struct fields to query parameters
func (p ScrapeParams) toMap() url.Values {
	return structToQueryParams(p)
}

// Generic function to convert a struct to URL query parameters
func structToQueryParams(data interface{}) url.Values {
	if data == nil {
		return url.Values{}
	}
	values := url.Values{}
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
			tagParts := strings.Split(tag, ",")
			if len(tagParts) > 1 && tagParts[1] == "omitempty" && value.IsZero() {
				continue
			}
			values.Set(tagParts[0], fmt.Sprintf("%v", value.Interface()))
		}
	}
	return values
}
