package client

// TODO: remove
import (
	"bytes"
	"fmt"
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
	// TODO: fill out
}

type TaggedImage struct {
	// TODO: fill out
}

const predictUrl string = "https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/versions/aa7f35c01e0642fda5cf400f543e7c40/outputs"

func (c Client) TagUrls(urls []string) error /* map[string]interface{} */ {
	var reqBody GeneralRequest
	urlCount := len(urls)
	inputs := make([]Inputs, urlCount)
	
	for i, url := range urls {
		inputs[i] = Inputs{ ImageData{ Image{ URL: url } } }
	}
	
	reqBody.Inputs = inputs
	jsonBody, err := json.Marshal(reqBody)

	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", predictUrl, bytes.NewReader(jsonBody))

	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "Key " + c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(resBody, &jsonResponse)

	if err != nil {
		return err
	}

	images := jsonResponse["outputs"].([]interface{})
	for _, image := range images {
		predictions := image.(map[string]interface{})["data"].(map[string]interface{})["concepts"].([]interface{})
		for _, prediction := range predictions {
			p := prediction.(map[string]interface{})
			fmt.Println(p["id"], p["name"])
		}
	}

	return nil
}
