package client

import (
	"encoding/json"
	"io/ioutil"
)

const config = "config/config.json"

// Client contains api key
type Client struct {
	APIKey string
}

// NewClient initializes a new client with the api key
func NewClient() (Client, error) {
	configFile, err := ioutil.ReadFile(config)
	var client Client
	
	if err != nil {
		return client, err
	}
	
	var config map[string]string
	err = json.Unmarshal(configFile, &config)

	if err != nil {
		return client, err
	}
	
	key := config["API_KEY"]
	client.APIKey = key

	return client, nil
}
