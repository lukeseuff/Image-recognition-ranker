package client

import (
	"container/heap"
	"fmt"
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

type ConceptHeap []Concept

func (h ConceptHeap) Len()          int     { return len(h) }
func (h ConceptHeap) Less(i, j int) bool    { return h[i].GetValue() < h[j].GetValue() }
func (h ConceptHeap) Swap(i, j int)         { h[i], h[j] = h[j], h[i] }
func (h ConceptHeap) Peek()         float64 { return h[0].GetValue() }

func Sort(h ConceptHeap) []Concept {
	sorted := make([]Concept, len(h))
	length := len(h)
	for i := 0; i < length; i++ {
		sorted[i] = h.Pop().(Concept)
		fmt.Printf("%v - %v\n", h, sorted[i].GetValue())
	}
	return sorted
}

func (h *ConceptHeap) Push(x interface{}) {
	*h = append(*h, x.(Concept))
}

func (h *ConceptHeap) Pop() interface {} {
	old := *h
	n := len(old)
	x := old[n - 1]
	*h = old[0 : n - 1]
	return x
}

func (c Client) TagURLs(urls []string) ([]TaggedImage, map[string][]Concept, error) {	
	conceptValues := make(map[string][]Concept)
	concepts := make(map[string]*ConceptHeap)
	urlCount := len(urls)

	taggedImages := make([]TaggedImage, urlCount)

	jsonRes, err := c.RequestPrediction(urls)

	if err != nil {
		return taggedImages, conceptValues, err
	}
	
	images := jsonRes["outputs"].([]interface{})
	
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
				if conceptList.Len() < 10 {
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

	for concept, conceptHeap := range concepts {
		conceptValues[concept] = Sort(*conceptHeap)
	}

	return taggedImages, conceptValues, nil
}
