package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ListLocations(pageURL *string) (ResShallowLocations, error) {

	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	if data, ok := c.cache.Get(url); ok {
		// test print
		fmt.Println("cache hit:", url)
		locationsRes := ResShallowLocations{}
		err := json.Unmarshal(data, &locationsRes)
		if err != nil {
			return ResShallowLocations{}, err
		}
		return locationsRes, nil
	}

	fmt.Println("cache miss:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ResShallowLocations{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return ResShallowLocations{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return ResShallowLocations{}, err
	}

	c.cache.Add(url, data)

	locationsRes := ResShallowLocations{}
	err = json.Unmarshal(data, &locationsRes)
	if err != nil {
		return ResShallowLocations{}, err
	}

	return locationsRes, nil
}
