package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) GetLocation(locationName string) (Location, error) {

	url := baseURL + "/location-area/" + locationName

	if data, ok := c.cache.Get(url); ok {

		locationsRes := Location{}
		err := json.Unmarshal(data, &locationsRes)
		if err != nil {
			return Location{}, err
		}
		return locationsRes, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Location{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Location{}, err
	}

	c.cache.Add(url, data)

	locationsRes := Location{}
	err = json.Unmarshal(data, &locationsRes)
	if err != nil {
		return Location{}, err
	}

	return locationsRes, nil
}
