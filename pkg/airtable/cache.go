package airtable

import (
	"log"
	"time"
)

type cachedResponse struct {
	retrieved time.Time
	response  string
}

type Cache struct {
	Responses map[string]cachedResponse
	Limit     time.Duration // cache for 1min
}

type Retriever func() string

func (c *Cache) Get(path string, retriever Retriever) string {
	if resp, prs := c.Responses[path]; prs {
		if time.Since(resp.retrieved) < c.Limit {
			log.Printf("  Cache hit %s\n", path)
			return resp.response
		}
	}

	resp := retriever()
	c.Responses[path] = cachedResponse{
		retrieved: time.Now(),
		response:  resp,
	}
	return resp
}

func (c *Cache) Purge() {
	c.Responses = make(map[string]cachedResponse)
}
