package ujeebu

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockSerpServer(response string, headers map[string]string, contentType string, statusCode int) (*httptest.Server, *Client) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range headers {
			w.Header().Add(k, v)
		}
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().SetBaseURL(mockServer.URL).SetTimeout(10 * time.Second),
	}

	return mockServer, client
}

func TestSerp_Success(t *testing.T) {
	mockResponse := `{"items":["Result 1","Result 2"]}`

	mockServer, client := setupMockSerpServer(mockResponse, map[string]string{
		"ujb-credits": "3",
	}, "application/json", http.StatusOK)
	defer mockServer.Close()

	params := SerpParams{
		Search:       "Go programming",
		SearchType:   "text",
		Lang:         "en",
		Location:     "us",
		Device:       "desktop",
		ResultsCount: 10,
		Page:         1,
	}

	serp, credits, err := client.Serp(params)
	require.NoError(t, err)
	assert.NotNil(t, serp)
	assert.Equal(t, 3, credits)
	resp := map[string]interface{}{}
	err = json.Unmarshal(serp, &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "items")
	assert.Contains(t, resp["items"], "Result 1")
}

func TestSerp_ErrorResponse(t *testing.T) {
	mockResponse := `{
		"status": "error",
		"error": "Invalid search query"
	}`

	mockServer, client := setupMockSerpServer(mockResponse, map[string]string{}, "application/json", http.StatusBadRequest)
	defer mockServer.Close()

	params := SerpParams{
		Search: "",
	}

	serp, credits, err := client.Serp(params)

	assert.Error(t, err)
	assert.Nil(t, serp)
	assert.Equal(t, 0, credits)
	assert.Contains(t, err.Error(), "Invalid search query")
}

func TestSerp_Timeout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "results": "{\"items\":[\"Delayed Result\"]}"}`))
	}))
	defer mockServer.Close()

	client := &Client{
		apiKey: "test_api_key",
		client: resty.New().SetBaseURL(mockServer.URL).SetTimeout(1 * time.Second),
	}

	params := SerpParams{
		Search: "Go programming",
	}

	serp, credits, err := client.Serp(params)

	assert.Error(t, err)
	assert.Nil(t, serp)
	assert.Equal(t, 0, credits)
}

func TestGoogleSearchResponse(t *testing.T) {
	mockResponse := `{
		"knowledge_graph": {
			"born": "July 10, 1856, Smiljan, Croatia",
			"died": "January 7, 1943 (age 86 years), The New Yorker A Wyndham Hotel, New York, NY",
			"education": "TU Graz (1875–1878), Gimnazija Karlovac (1870–1873)",
			"height": "6′ 2″",
			"parents": "Milutin Tesla, Ðuka Tesla",
			"siblings": "Dane Tesla, Angelina Tesla, Milka Tesla, Marica Kosanović",
			"title": "Nikola Tesla",
			"type": "Engineer and futurist"
		},
		"metadata": {
			"google_url": "https://www.google.com/search?gl=US&hl=en&num=10&q=Nikola+Tesla&sei=defQZ8riGZOt5NoPk7_S4AU",
			"number_of_results": 36800000,
			"query_displayed": "Nikola Tesla",
			"results_time": "0.29 seconds"
		},
		"organic_results": [
			{
				"cite": "https://en.wikipedia.org › wiki › Nikola_Tesla",
				"link": "https://en.wikipedia.org/wiki/Nikola_Tesla",
				"position": 1,
				"site_name": "Wikipedia",
				"title": "Nikola Tesla"
			},
			{
				"cite": "https://www.britannica.com › ... › Matter & Energy",
				"link": "https://www.britannica.com/biography/Nikola-Tesla",
				"position": 2,
				"site_name": "Britannica",
				"title": "Nikola Tesla | Biography, Facts, & Inventions"
			}
		]
	}`

	mockServer, client := setupMockSerpServer(mockResponse, map[string]string{
		"ujb-credits": "5",
	}, "application/json", http.StatusOK)
	defer mockServer.Close()

	params := SerpParams{
		Search:       "Nikola Tesla",
		SearchType:   "text",
		Lang:         "en",
		Location:     "us",
		Device:       "desktop",
		ResultsCount: 10,
		Page:         1,
	}

	result, credits, err := client.GoogleSearch(params)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 5, credits)

	// Validate Knowledge Graph
	assert.Equal(t, "Nikola Tesla", result.KnowledgeGraph.Title)
	assert.Equal(t, "Engineer and futurist", result.KnowledgeGraph.Type)
	assert.Contains(t, result.KnowledgeGraph.Parents, "Milutin Tesla")

	// Validate Metadata
	assert.Equal(t, "https://www.google.com/search?gl=US&hl=en&num=10&q=Nikola+Tesla&sei=defQZ8riGZOt5NoPk7_S4AU", result.Metadata.GoogleUrl)
	assert.Equal(t, 36800000, result.Metadata.NumberOfResults)

	// Validate Organic Results
	require.Len(t, result.OrganicResults, 2)
	assert.Equal(t, "Wikipedia", result.OrganicResults[0].SiteName)
	assert.Equal(t, "Nikola Tesla", result.OrganicResults[0].Title)
	assert.Equal(t, 1, result.OrganicResults[0].Position)

	assert.Equal(t, "Britannica", result.OrganicResults[1].SiteName)
	assert.Contains(t, result.OrganicResults[1].Link, "britannica.com")
}

func TestGoogleNewsResponse(t *testing.T) {
	mockResponse := `{
		"metadata": {
			"google_url": "https://www.google.com/search?gl=US&hl=en&num=10&q=Donald+Trump&tbm=nws",
			"number_of_results": 62800000,
			"query_displayed": "Donald Trump",
			"results_time": "0.25 seconds"
		},
		"news": [
			{
				"date": "LIVE2 hours ago",
				"description": "The House has passed the spending measure to fund the government through September 30. Meanwhile, President Volodymyr Zelensky said Ukraine...",
				"favicon": "data:image/png;base64,iVBOR...",
				"image": "data:image/jpeg;base64,/9j/4AA...",
				"link": "https://www.cnn.com/politics/live-news/trump-administration-presidency-ukraine-03-11-2025/index.html",
				"position": 1,
				"siteName": "CNN",
				"title": "Live updates: Trump tariff threats; Government funding bill House vote; US-Ukraine talks"
			},
			{
				"date": "30 minutes ago",
				"description": "President Donald Trump's threat to double his planned tariffs on steel and aluminum from 25% to 50% for Canada has led the provincial...",
				"favicon": "data:image/png;base64,iVBOR...",
				"image": "data:image/jpeg;base64,/9j/4AA...",
				"link": "https://apnews.com/article/trump-economy-tariffs-stock-musk-business-8a5f28d9bb16e0b8a924d99ead0907fa",
				"position": 2,
				"siteName": "AP News",
				"title": "Trump halts doubling of tariffs on Canadian metals after Ontario suspends electricity price hikes"
			}
		],
		"pagination": {
			"google": {
				"current": "https://google.com/search?gl=US&hl=en&num=10&q=Donald+Trump&tbm=nws&",
				"next": "https://google.com/search?gl=US&hl=en&num=10&q=Donald+Trump&start=20&tbm=&"
			},
			"api": {
				"current": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=1&results_count=10&search=Donald+Trump&",
				"next": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=2&results_count=10&search=Donald+Trump&"
			}
		}
	}`

	mockServer, client := setupMockSerpServer(mockResponse, map[string]string{
		"ujb-credits": "7",
	}, "application/json", http.StatusOK)
	defer mockServer.Close()

	params := SerpParams{
		Search:       "Donald Trump",
		SearchType:   "news",
		Lang:         "en",
		Location:     "us",
		Device:       "desktop",
		ResultsCount: 10,
		Page:         1,
	}

	result, credits, err := client.GoogleNewsSearch(params)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 7, credits)

	// Validate Metadata
	assert.Equal(t, "https://www.google.com/search?gl=US&hl=en&num=10&q=Donald+Trump&tbm=nws", result.Metadata.GoogleUrl)
	assert.Equal(t, 62800000, result.Metadata.NumberOfResults)
	assert.Equal(t, "Donald Trump", result.Metadata.QueryDisplayed)

	// Validate News
	require.Len(t, result.News, 2)
	assert.Equal(t, "CNN", result.News[0].SiteName)
	assert.Equal(t, "Live updates: Trump tariff threats; Government funding bill House vote; US-Ukraine talks", result.News[0].Title)
	assert.Equal(t, 1, result.News[0].Position)

	assert.Equal(t, "AP News", result.News[1].SiteName)
	assert.Equal(t, "Trump halts doubling of tariffs on Canadian metals after Ontario suspends electricity price hikes", result.News[1].Title)
	assert.Equal(t, 2, result.News[1].Position)

	// Validate Pagination
	assert.Equal(t, "https://google.com/search?gl=US&hl=en&num=10&q=Donald+Trump&tbm=nws&", result.Pagination.Google.Current)
	assert.Equal(t, "https://google.com/search?gl=US&hl=en&num=10&q=Donald+Trump&start=20&tbm=&", result.Pagination.Google.Next)
	assert.Equal(t, "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=1&results_count=10&search=Donald+Trump&", result.Pagination.Api.Current)
	assert.Equal(t, "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=2&results_count=10&search=Donald+Trump&", result.Pagination.Api.Next)
}

func TestGoogleVideoResponse(t *testing.T) {
	mockResponse := `{
		"metadata": {
			"google_url": "https://www.google.com/search?gl=US&hl=en&num=10&q=Bitcoin&udm=7",
			"query_displayed": "Bitcoin"
		},
		"pagination": {
			"google": {
				"current": "https://google.com/search?gl=US&hl=en&num=10&q=Bitcoin&tbm=vid&",
				"next": "https://google.com/search?gl=US&hl=en&num=10&q=Bitcoin&start=20&tbm=&",
				"other_pages": {
					"3": "https://google.com/search?gl=US&hl=en&num=10&q=Bitcoin&start=30&tbm=&",
					"4": "https://google.com/search?gl=US&hl=en&num=10&q=Bitcoin&start=40&tbm=&",
					"5": "https://google.com/search?gl=US&hl=en&num=10&q=Bitcoin&start=50&tbm=&",
					"6": "https://google.com/search?gl=US&hl=en&num=10&q=Bitcoin&start=60&tbm=&",
					"7": "https://google.com/search?gl=US&hl=en&num=10&q=Bitcoin&start=70&tbm=&",
					"8": "https://google.com/search?gl=US&hl=en&num=10&q=Bitcoin&start=80&tbm=&"
				}
			},
			"api": {
				"current": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=1&results_count=10&search=Bitcoin&",
				"next": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=2&results_count=10&search=Bitcoin&",
				"other_pages": {
					"3": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=3&results_count=10&search=Bitcoin&",
					"4": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=4&results_count=10&search=Bitcoin&",
					"5": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=5&results_count=10&search=Bitcoin&",
					"6": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=6&results_count=10&search=Bitcoin&",
					"7": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=7&results_count=10&search=Bitcoin&",
					"8": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=8&results_count=10&search=Bitcoin&"
				}
			}
		},
		"videos": [
			{
				"author": "7 hours ago",
				"date": "7 hours ago",
				"description": "On today's episode of CNBC Crypto World, bitcoin jumps back above $80000 after a Monday sell-off. Plus, senators reintroduce the GENIUS Act...",
				"position": 1,
				"provider": "CNBC",
				"site": "www.cnbc.com",
				"summary": "Bitcoin rebounds after falling to its lowest level since November on CNBC. Play on CNBC. 10:02. 7 hours ago",
				"title": "Bitcoin rebounds after falling to its lowest level since November",
				"url": "https://www.cnbc.com/video/2025/03/11/bitcoin-rebounds-after-lowest-level-since-november-crypto-world.html"
			},
			{
				"author": "CNBC Television",
				"date": "6 hours ago",
				"description": "On today's episode of CNBC Crypto World, bitcoin jumps back above $80000 after a Monday sell-off. Plus, senators reintroduce the GENIUS Act...",
				"position": 2,
				"provider": "YouTube",
				"site": "www.youtube.com",
				"summary": "Bitcoin rebounds after falling to its lowest level since ... by CNBC Television on YouTube. Play on YouTube. 10:03. 6 hours ago",
				"title": "Bitcoin rebounds after falling to its lowest level since ...",
				"url": "https://www.youtube.com/watch?v=ts0VTnXMs9I"
			}
		]
	}`

	mockServer, client := setupMockSerpServer(mockResponse, map[string]string{
		"ujb-credits": "12",
	}, "application/json", http.StatusOK)
	defer mockServer.Close()

	params := SerpParams{
		Search:       "Bitcoin",
		SearchType:   "video",
		Lang:         "en",
		Location:     "us",
		Device:       "desktop",
		ResultsCount: 10,
		Page:         1,
	}

	result, credits, err := client.GoogleVideoSearch(params)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 12, credits)

	// Validate Metadata
	assert.Equal(t, "https://www.google.com/search?gl=US&hl=en&num=10&q=Bitcoin&udm=7", result.Metadata.GoogleUrl)
	assert.Equal(t, "Bitcoin", result.Metadata.QueryDisplayed)

	// Validate Videos
	require.Len(t, result.Videos, 2)
	assert.Equal(t, "CNBC", result.Videos[0].Provider)
	assert.Equal(t, "Bitcoin rebounds after falling to its lowest level since November", result.Videos[0].Title)
	assert.Equal(t, 1, result.Videos[0].Position)

	assert.Equal(t, "YouTube", result.Videos[1].Provider)
	assert.Equal(t, "Bitcoin rebounds after falling to its lowest level since ...", result.Videos[1].Title)
	assert.Equal(t, 2, result.Videos[1].Position)

	// Validate Pagination
	assert.Equal(t, "https://google.com/search?gl=US&hl=en&num=10&q=Bitcoin&tbm=vid&", result.Pagination.Google.Current)
	assert.Equal(t, "https://google.com/search?gl=US&hl=en&num=10&q=Bitcoin&start=20&tbm=&", result.Pagination.Google.Next)
	assert.Equal(t, "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=1&results_count=10&search=Bitcoin&", result.Pagination.Api.Current)
	assert.Equal(t, "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=2&results_count=10&search=Bitcoin&", result.Pagination.Api.Next)
}

func TestGoogleImageResponse(t *testing.T) {
	mockResponse := `{
		"images": [
			{
				"google_thumbnail": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRkRapDtN6JSis1bWCnMbqn3pmIEDeDY9t8tg\u0026s",
				"height": 2000,
				"image": "https://upload.wikimedia.org/wikipedia/commons/e/e4/Latte_and_dark_coffee.jpg",
				"link": "https://en.wikipedia.org/wiki/Coffee",
				"position": 1,
				"source": "Wikipedia",
				"title": "Coffee - Wikipedia",
				"width": 3200
			},
			{
				"google_thumbnail": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTsmD2K-wg4D2csngncAOMlxa1y5g_yToyicw\u0026s",
				"height": 405,
				"image": "https://somedayilllearn.com/wp-content/uploads/2020/05/cup-of-black-coffee-scaled-720x405.jpeg",
				"link": "https://somedayilllearn.com/how-to-make-black-coffee/",
				"position": 2,
				"source": "Someday I'll Learn",
				"title": "How to Make Black Coffee that Tastes Good",
				"width": 720
			}
		],
		"metadata": {
			"google_url": "https://www.google.com/search?gl=US\u0026hl=en\u0026num=10\u0026q=Coffee\u0026udm=2",
			"query_displayed": "Coffee"
		},
		"pagination": {
			"google": {
				"current": "https://google.com/search?gl=US\u0026hl=en\u0026num=10\u0026q=Coffee\u0026udm=2\u0026",
				"next": "https://google.com/search?gl=US\u0026hl=en\u0026num=10\u0026q=Coffee\u0026start=20\u0026tbm=\u0026"
			},
			"api": {
				"current": "https://api.ujeebu.com/serp?device=desktop\u0026lang=en\u0026location=US\u0026page=1\u0026results_count=10\u0026search=Coffee\u0026",
				"next": "https://api.ujeebu.com/serp?device=desktop\u0026lang=en\u0026location=US\u0026page=2\u0026results_count=10\u0026search=Coffee\u0026"
			}
		}
	}`

	mockServer, client := setupMockSerpServer(mockResponse, map[string]string{
		"ujb-credits": "15",
	}, "application/json", http.StatusOK)
	defer mockServer.Close()

	params := SerpParams{
		Search:       "Coffee",
		SearchType:   "image",
		Lang:         "en",
		Location:     "us",
		Device:       "desktop",
		ResultsCount: 10,
		Page:         1,
	}

	result, credits, err := client.GoogleImageSearch(params)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 15, credits)

	// Validate Metadata
	assert.Equal(t, "https://www.google.com/search?gl=US&hl=en&num=10&q=Coffee&udm=2", result.Metadata.GoogleUrl)
	assert.Equal(t, "Coffee", result.Metadata.QueryDisplayed)

	// Validate Images
	require.Len(t, result.Images, 2)
	assert.Equal(t, "Wikipedia", result.Images[0].Source)
	assert.Equal(t, "Coffee - Wikipedia", result.Images[0].Title)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Coffee", result.Images[0].Link)
	assert.Equal(t, 1, result.Images[0].Position)
	assert.Equal(t, "https://upload.wikimedia.org/wikipedia/commons/e/e4/Latte_and_dark_coffee.jpg", result.Images[0].Image)
	assert.Equal(t, 2000, result.Images[0].Height)
	assert.Equal(t, 3200, result.Images[0].Width)

	assert.Equal(t, "Someday I'll Learn", result.Images[1].Source)
	assert.Equal(t, "How to Make Black Coffee that Tastes Good", result.Images[1].Title)
	assert.Equal(t, "https://somedayilllearn.com/how-to-make-black-coffee/", result.Images[1].Link)
	assert.Equal(t, 2, result.Images[1].Position)
	assert.Equal(t, "https://somedayilllearn.com/wp-content/uploads/2020/05/cup-of-black-coffee-scaled-720x405.jpeg", result.Images[1].Image)
	assert.Equal(t, 405, result.Images[1].Height)
	assert.Equal(t, 720, result.Images[1].Width)

	// Validate Pagination
	assert.Equal(t, "https://google.com/search?gl=US&hl=en&num=10&q=Coffee&udm=2&", result.Pagination.Google.Current)
	assert.Equal(t, "https://google.com/search?gl=US&hl=en&num=10&q=Coffee&start=20&tbm=&", result.Pagination.Google.Next)
	assert.Equal(t, "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=1&results_count=10&search=Coffee&", result.Pagination.Api.Current)
	assert.Equal(t, "https://api.ujeebu.com/serp?device=desktop&lang=en&location=US&page=2&results_count=10&search=Coffee&", result.Pagination.Api.Next)
}

func TestGoogleMapSearchResponse(t *testing.T) {
	mockResponse := `{
		"maps_results": [
			{
				"address": "Brampton, ON",
				"category": "\u00b7",
				"cid": "1322681635158113553",
				"opening_hours": null,
				"position": 1,
				"rating": 4.7,
				"reviews": 370,
				"title": "La Pergola Ristorante"
			},
			{
				"address": "Burnaby, BC",
				"category": "Italian",
				"cid": "5098284575974088326",
				"opening_hours": null,
				"position": 2,
				"rating": 4.6,
				"reviews": 598,
				"title": "Trattoria by Italian Kitchen"
			}
		],
		"metadata": {
			"google_url": "https://www.google.com/search?gl=ca&hl=en&num=10&q=Italian+restaurant&tbm=lcl",
			"query_displayed": "Italian restaurant"
		},
		"pagination": {
			"google": {
				"current": "https://google.com/search?gl=ca&hl=en&num=10&q=Italian+restaurant&tbm=lcl&",
				"next": "https://google.com/search?gl=ca&hl=en&num=10&q=Italian+restaurant&start=20&tbm=&"
			},
			"api": {
				"current": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=ca&page=1&results_count=10&search=Italian+restaurant&",
				"next": "https://api.ujeebu.com/serp?device=desktop&lang=en&location=ca&page=2&results_count=10&search=Italian+restaurant&"
			}
		}
	}`

	mockServer, client := setupMockSerpServer(mockResponse, map[string]string{
		"ujb-credits": "20",
	}, "application/json", http.StatusOK)
	defer mockServer.Close()

	params := SerpParams{
		Search:       "Italian restaurant",
		SearchType:   "map",
		Lang:         "en",
		Location:     "ca",
		Device:       "desktop",
		ResultsCount: 10,
		Page:         1,
	}

	result, credits, err := client.GoogleMapSearch(params)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 20, credits)

	// Validate Metadata
	assert.Equal(t, "https://www.google.com/search?gl=ca&hl=en&num=10&q=Italian+restaurant&tbm=lcl", result.Metadata.GoogleUrl)
	assert.Equal(t, "Italian restaurant", result.Metadata.QueryDisplayed)

	// Validate Maps Results
	require.Len(t, result.Maps, 2)
	assert.Equal(t, "La Pergola Ristorante", result.Maps[0].Title)
	assert.Equal(t, "Brampton, ON", result.Maps[0].Address)
	assert.Equal(t, 4.7, result.Maps[0].Rating)
	assert.Equal(t, 370, result.Maps[0].Reviews)
	assert.Equal(t, 1, result.Maps[0].Position)
	assert.Equal(t, "1322681635158113553", result.Maps[0].Cid)

	assert.Equal(t, "Trattoria by Italian Kitchen", result.Maps[1].Title)
	assert.Equal(t, "Burnaby, BC", result.Maps[1].Address)
	assert.Equal(t, 4.6, result.Maps[1].Rating)
	assert.Equal(t, 598, result.Maps[1].Reviews)
	assert.Equal(t, 2, result.Maps[1].Position)
	assert.Equal(t, "5098284575974088326", result.Maps[1].Cid)

	// Validate Pagination
	assert.Equal(t, "https://google.com/search?gl=ca&hl=en&num=10&q=Italian+restaurant&tbm=lcl&", result.Pagination.Google.Current)
	assert.Equal(t, "https://google.com/search?gl=ca&hl=en&num=10&q=Italian+restaurant&start=20&tbm=&", result.Pagination.Google.Next)
	assert.Equal(t, "https://api.ujeebu.com/serp?device=desktop&lang=en&location=ca&page=1&results_count=10&search=Italian+restaurant&", result.Pagination.Api.Current)
	assert.Equal(t, "https://api.ujeebu.com/serp?device=desktop&lang=en&location=ca&page=2&results_count=10&search=Italian+restaurant&", result.Pagination.Api.Next)
}
