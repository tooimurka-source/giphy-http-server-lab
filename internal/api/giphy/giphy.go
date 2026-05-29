package giphy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	Data []struct {
		Images struct {
			Original struct {
				URL string `json:"url"`
			} `json:"original"`
		} `json:"images"`
	} `json:"data"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Search(query string) (string, error) {
	endpoint, err := url.Parse("https://api.giphy.com/v1/gifs/search")
	if err != nil {
		return "", err
	}

	params := endpoint.Query()
	params.Set("api_key", c.apiKey)
	params.Set("q", query)
	params.Set("limit", "1")
	endpoint.RawQuery = params.Encode()

	resp, err := c.httpClient.Get(endpoint.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("giphy returned status %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Data) == 0 {
		return "", fmt.Errorf("no gifs found")
	}

	return result.Data[0].Images.Original.URL, nil
}

func GetGif(apiKey, query string) (string, error) {
	return NewClient(apiKey).Search(query)
}
