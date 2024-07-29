package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
	"time"
)

type LocationAreaDetails struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance   int `json:"chance"`
				MaxLevel int `json:"max_level"`
				MinLevel int `json:"min_level"`
				Method   struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Client struct {
	BaseURL    string
	cache      *pokecache.Cache
	httpClient *http.Client
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

func NewClient(cacheInterval time.Duration) *Client {
	return &Client{
		BaseURL:    "https://pokeapi.co/api/v2/",
		cache:      pokecache.NewCache(cacheInterval),
		httpClient: &http.Client{},
	}
}

func (c *Client) GetLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	endpoint := fmt.Sprintf("%s/location-area", c.BaseURL)
	if pageURL != nil {
		endpoint = *pageURL
	}
	cache, ok := c.cache.Get(endpoint)
	if ok {
		var locationResp LocationAreaResponse
		err := json.Unmarshal(cache, &locationResp)
		if err == nil {
			return locationResp, nil
		}
	}
	res, err := c.httpClient.Get(endpoint)
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

	c.cache.Add(endpoint, body)

	var locationResp LocationAreaResponse
	err = json.Unmarshal(body, &locationResp)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return locationResp, nil
}

func (c *Client) GetLocationArea(name string) (LocationAreaDetails, error) {
	endpoint := fmt.Sprintf("%s/location-area/%s", c.BaseURL, name)

	if cachedData, ok := c.cache.Get(endpoint); ok {
		var locationArea LocationAreaDetails
		err := json.Unmarshal(cachedData, &locationArea)
		if err == nil {
			return locationArea, nil
		}
	}
	res, err := c.httpClient.Get(endpoint)
	if err != nil {
		return LocationAreaDetails{}, fmt.Errorf("error making GET request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaDetails{}, fmt.Errorf("error reading response body: %w", err)
	}

	if res.StatusCode > 299 {
		return LocationAreaDetails{}, fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}

	var locationArea LocationAreaDetails
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return LocationAreaDetails{}, fmt.Errorf("error unmarshalling response body: %w", err)
	}
	c.cache.Add(endpoint, body)
	return locationArea, nil
}
