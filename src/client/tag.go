package client

// TODO: remove
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

type GeneralRequest struct {
	Inputs []Inputs `json:"inputs"`
}

type Tag struct {
	Concept     string
	Probability float64
}

type TaggedImage struct {
	URL  string
	Tags []Tag
}

const predictURL string = "https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/versions/aa7f35c01e0642fda5cf400f543e7c40/outputs"

func (c Client) TagURLs(urls []string) ([]TaggedImage, error) {
	var reqBody GeneralRequest
	urlCount := len(urls)
	inputs := make([]Inputs, urlCount)
	taggedImages := make([]TaggedImage, urlCount)
	
	for i, url := range urls {
		inputs[i] = Inputs{ ImageData{ Image{ URL: url } } }
	}
	
	reqBody.Inputs = inputs
	jsonBody, err := json.Marshal(reqBody)

	if err != nil {
		return taggedImages, err
	}
	
	req, err := http.NewRequest("POST", predictURL, bytes.NewReader(jsonBody))

	if err != nil {
		return taggedImages, err
	}
	
	req.Header.Set("Authorization", "Key " + c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return taggedImages, err
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return taggedImages, err
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(resBody, &jsonResponse)

	if err != nil {
		return taggedImages, err
	}

	images := jsonResponse["outputs"].([]interface{})
	for imageIndex, image := range images {
		var taggedImage TaggedImage
		imageURL := image.(map[string]interface{})["input"].(map[string]interface{})["data"].(map[string]interface{})["image"].(map[string]interface{})["url"].(string)
		predictions := image.(map[string]interface{})["data"].(map[string]interface{})["concepts"].([]interface{})
		taggedImage.URL = imageURL
		for _, prediction := range predictions {
			p := prediction.(map[string]interface{})
			tag := Tag{ Concept: p["name"].(string), Probability: p["value"].(float64) }
			// TODO: allocate size upfront
			taggedImage.Tags = append(taggedImage.Tags, tag)
		}
		taggedImages[imageIndex] = taggedImage
	}

	return taggedImages, nil
}
