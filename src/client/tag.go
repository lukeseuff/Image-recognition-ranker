package client

// TODO: remove
import (
	"bytes"
	"container/heap"
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

type TaggedImage struct {
	URL     string
	Concept map[string]float64
}

type Concept struct {
	Name  string
	Image *TaggedImage
}

type ConceptHeap []Concept

func (h ConceptHeap) Len()          int     { return len(h) }
func (h ConceptHeap) Less(i, j int) bool    { return h[i].Image.Concept[h[i].Name] < h[j].Image.Concept[h[j].Name] }
func (h ConceptHeap) Swap(i, j int)         { h[i], h[j] = h[j], h[i] }
func (h ConceptHeap) Peek()         float64 { return h[0].Image.Concept[h[0].Name] }

func (h *ConceptHeap) Push(x interface{}) {
	*h = append(*h, x.(Concept))
}

func (h *ConceptHeap) Pop() interface {} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

const predictURL string = "https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/versions/aa7f35c01e0642fda5cf400f543e7c40/outputs"

func (c Client) TagURLs(urls []string) ([]TaggedImage, map[string]*ConceptHeap, error) {
	var reqBody GeneralRequest
	concepts := make(map[string]*ConceptHeap)
	urlCount := len(urls)
	inputs := make([]Inputs, urlCount)
	taggedImages := make([]TaggedImage, urlCount)
	
	for i, url := range urls {
		inputs[i] = Inputs{ ImageData{ Image{ URL: url } } }
	}
	
	reqBody.Inputs = inputs
	jsonBody, err := json.Marshal(reqBody)

	if err != nil {
		return taggedImages, concepts, err
	}
	
	req, err := http.NewRequest("POST", predictURL, bytes.NewReader(jsonBody))

	if err != nil {
		return taggedImages, concepts, err
	}
	
	req.Header.Set("Authorization", "Key " + c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return taggedImages, concepts, err
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return taggedImages, concepts, err
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(resBody, &jsonResponse)

	if err != nil {
		return taggedImages, concepts, err
	}

	images := jsonResponse["outputs"].([]interface{})
	for imageIndex, image := range images {
		imageURL := image.(map[string]interface{})["input"].(map[string]interface{})["data"].(map[string]interface{})["image"].(map[string]interface{})["url"].(string)
		predictions := image.(map[string]interface{})["data"].(map[string]interface{})["concepts"].([]interface{})
		taggedImage := TaggedImage{ URL: imageURL, Concept: make(map[string]float64) }
		taggedImages[imageIndex] = taggedImage
		for _, prediction := range predictions {
			p := prediction.(map[string]interface{})
			taggedImages[imageIndex].Concept[p["name"].(string)] = p["value"].(float64)
			concept := Concept{ Name: p["name"].(string), Image: &taggedImages[imageIndex] }
		
			conceptList, ok := concepts[p["name"].(string)]
			if ok {
				if conceptList.Len() <= 10 {
					heap.Push(conceptList, concept)
				} else if concept.Image.Concept[concept.Name] > conceptList.Peek() {
					conceptList.Pop()
					heap.Push(conceptList, concept)
				}
			} else {
				newConcept := &ConceptHeap{ Concept{ Name: p["name"].(string), Image: &taggedImages[imageIndex] } }
				heap.Init(newConcept)
				concepts[p["name"].(string)] = newConcept
			}
		}
	}

	return taggedImages, concepts, nil
}
