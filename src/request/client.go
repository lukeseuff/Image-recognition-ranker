package request

import (
	"util"
)

const config = "config/config.json"

// Client contains api key
type Client struct {
	APIKey string
}

// NewClient initializes a new client with the api key
func NewClient() (Client, error) {
	var client Client
	config := util.GetConfig()
	key := config["API_KEY"].(string)
	client.APIKey = key

	return client, nil
}
