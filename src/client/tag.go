package client

import (
	"container/heap"
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
		sorted[length - i - 1] = heap.Pop(&h).(Concept)
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

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func (c Client) TagURLs(urls []string) (map[string][]Concept, error) {	
	conceptValues := make(map[string][]Concept)
	concepts := make(map[string]*ConceptHeap)
	urlCount := len(urls)

	batches := (len(urls) - 1)/128 + 1
	
	for batch := 0; batch < batches; batch++  {
		jsonRes, err := c.RequestPrediction(urls[batch*128 : min((batch + 1)*128, urlCount)])
		if err != nil {
			return conceptValues, err
		}
		images := jsonRes["outputs"].([]interface{})

		for _, image := range images {
			imageURL := image.(map[string]interface{})["input"].(map[string]interface{})["data"].(map[string]interface{})["image"].(map[string]interface{})["url"].(string)
			predictions := image.(map[string]interface{})["data"].(map[string]interface{})["concepts"].([]interface{})
			taggedImage := TaggedImage{ URL: imageURL, Concept: make(map[string]float64) }
			for _, prediction := range predictions {
				p := prediction.(map[string]interface{})
				concept := Concept{ Name: p["name"].(string), Image: &taggedImage }
				taggedImage.Concept[concept.Name] = p["value"].(float64)
				
				conceptList, ok := concepts[p["name"].(string)]
				if ok {
					if conceptList.Len() < 10 {
						heap.Push(conceptList, concept)
					} else if concept.GetValue() > conceptList.Peek() {
						heap.Pop(conceptList)
						heap.Push(conceptList, concept)
					}
				} else {
					newConcept := &ConceptHeap{ Concept{ Name: p["name"].(string), Image: &taggedImage } }
					heap.Init(newConcept)
					concepts[p["name"].(string)] = newConcept
				}
			}
		}
	}

	for concept, conceptHeap := range concepts {
		conceptValues[concept] = Sort(*conceptHeap)
	}

	return conceptValues, nil
}
