package client

import (
	"errors"
	"encoding/json"
	"io/ioutil"
)

const config = "../config/config.json"

type Client struct {
	APIKey string
}

func NewClient() (Client, error) {
	configFile, err := ioutil.ReadFile(config)
	var client Client
	
	if err != nil {
		return client, err
	}
	
	var config map[string]interface{}
	err = json.Unmarshal(configFile, &config)

	if err != nil {
		return client, err
	}
	
	key, ok := config["API_KEY"].(string)

	if !ok {
		return client, errors.New("API Key is not a string")
	}

	client.APIKey = key

	return client, nil
}
