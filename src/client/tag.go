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

func (c Client) TagUrl(url string) error {
	var reqBody GeneralRequest
	reqBody.Inputs = append(reqBody.Inputs, Inputs{ ImageData{ Image{ URL: url } } })
	jsonBody, err := json.Marshal(reqBody)

	if err != nil {
		return err
	}

	predictUrl := "https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/versions/aa7f35c01e0642fda5cf400f543e7c40/outputs"
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

	jsonResponse := new(map[string]interface{})
	err = json.Unmarshal(resBody, jsonResponse)

	if err != nil {
		return err
	}

	// TODO: remove
	fmt.Printf("%+v\n", jsonResponse)
	
	return nil
}

// func (Client) TagUrls([]string urls) map[string]interface{} {
	
// }
