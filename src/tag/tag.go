package tag

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)


type TaggedImage struct {
	URL     string
	Concept map[string]float64
}

type Concept struct {
	Name  string
	Image *TaggedImage
}

func (c Concept) GetValue() float64 {
	return c.Image.Concept[c.Name]
}

func addConcepts(img *TaggedImage, predictions []interface{}) {
	for _, concept := range predictions {
		c := concept.(map[string]interface{})
		img.Concept[c["name"].(string)] = c["value"].(float64)
	}
}


func TagImages(responses []*http.Response) []*TaggedImage {
	var jsonResponse map[string]interface{}
	var taggedImages []*TaggedImage

	for _, response := range responses {
		responseBody, _ := ioutil.ReadAll(response.Body)
		_ = json.Unmarshal(responseBody, &jsonResponse)
		images := jsonResponse["outputs"].([]interface{})
		for _, image := range images {
			url := image.(map[string]interface{})["input"].
				(map[string]interface{})["data"].
				(map[string]interface{})["image"].
				(map[string]interface{})["url"].
				(string)
			predictions := image.(map[string]interface{})["data"].
				(map[string]interface{})["concepts"].
				([]interface{})
			taggedImage := &TaggedImage{ URL: url, Concept: make(map[string]float64, len(predictions)) }
			addConcepts(taggedImage, predictions)
			taggedImages = append(taggedImages, taggedImage)
		}
	}

	return taggedImages
}
