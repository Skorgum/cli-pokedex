package pokeapi

import (
	"net/http"
	"time"

	"github.com/Skorgum/cli-pokedex/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	cache      pokecache.Cache
}

func NewClient(timeout time.Duration) Client {
	const defaultCacheInterval = 5 * time.Second
	return Client{
		cache: *pokecache.NewCache(defaultCacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}
