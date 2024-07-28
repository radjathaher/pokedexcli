package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

func NewClient() *Client {
	return &Client{BaseURL: "https://pokeapi.co/api/v2/"}
}

func (c *Client) GetLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	endpoint := fmt.Sprintf("%s/location-area", c.BaseURL)
	if pageURL != nil {
		endpoint = *pageURL
	}
	res, err := http.Get(endpoint)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error making GET request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	if res.StatusCode > 299 {
		return LocationAreaResponse{}, fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}
	var locationResp LocationAreaResponse
	err = json.Unmarshal(body, &locationResp)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return locationResp, nil
}
