package request

import (
	"bytes"
	"net/http"
	"encoding/json"
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

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}


func (c Client) BatchPrediction(urls []string, batchSize int) []*http.Response {
	batches := (len(urls) - 1)/batchSize + 1
	responses := make([]*http.Response, batches)

	for batch := 0; batch < batches; batch++ {
		beginSlice := batch*batchSize
		endSlice := min((batch + 1)*batchSize, len(urls))
		responses[batch] = c.RequestPrediction(urls[beginSlice:endSlice])
	}

	return responses
}

func (c Client) RequestPrediction(urls []string) *http.Response {
	var requestBody InputRequest
	inputs := make([]Inputs, len(urls))

	for i, url := range urls {
		inputs[i] = Inputs{ ImageData{ Image{ URL: url } } }
	}

	requestBody.Inputs = inputs
	jsonRequest, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", predictURL, bytes.NewReader(jsonRequest))
	req.Header.Set("Authorization", "Key " + c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, _ := client.Do(req)

	return response
}
