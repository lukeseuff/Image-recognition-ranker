package client

import (
	"bytes"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type Image struct {
	URL string `json:"url"`
}

type ImageData struct {
	Image `json:"image"`
}

type Inputs struct {
	Data ImageData `json:"data"`
}

type InputRequest struct {
	Inputs []Inputs `json:"inputs"`
}

const predictURL string = "https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/versions/aa7f35c01e0642fda5cf400f543e7c40/outputs"

func (c Client) RequestPrediction(urls []string) (map[string]interface{}, error) {
	var body InputRequest
	var jsonRes map[string]interface{}
	
	inputs := make([]Inputs, len(urls))

	for i, url := range urls {
		inputs[i] = Inputs{ ImageData{ Image{ URL: url } } }
	}

	body.Inputs = inputs
	jsonBody, err := json.Marshal(body)

	if err != nil {
		return jsonRes, err
	}
	
	req, err := http.NewRequest("POST", predictURL, bytes.NewReader(jsonBody))

	if err != nil {
		return jsonRes, err
	}
	
	req.Header.Set("Authorization", "Key " + c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return jsonRes, err
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return jsonRes, err
	}

	err = json.Unmarshal(resBody, &jsonRes)

	if err != nil {
		return jsonRes, err
	}

	return jsonRes, nil
}
