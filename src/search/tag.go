package search

import "container/heap"

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

const batchSize = 128

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

func addConcepts(img *TaggedImage, predictions []interface{}) {
	for _, concept := range predictions {
		c := concept.(map[string]interface{})
		img.Concept[c["name"].(string)] = c["value"].(float64)
	}
}

func addImages(res []interface{}) []*TaggedImage {
	tagged := make([]*TaggedImage, 0, len(res))
	for _, r := range res {
		imageURL := r.(map[string]interface{})["input"].(map[string]interface{})["data"].(map[string]interface{})["image"].(map[string]interface{})["url"].(string)
		predictions := r.(map[string]interface{})["data"].(map[string]interface{})["concepts"].([]interface{})
		taggedImage := TaggedImage {
			URL: imageURL,
			Concept: make(map[string]float64),
		}
		tagged = append(tagged, &taggedImage)
		addConcepts(&taggedImage, predictions)
	}
	return tagged
}

func (c Client) TagImages(urls []string) ([]*TaggedImage, error) {
	tagged := make([]*TaggedImage, 0, len(urls))
	batches := (len(urls) - 1)/batchSize + 1
	
	for b := 0; b < batches; b++  {
		jsonRes, err := c.RequestPrediction(urls[b*128 : min((b + 1)*128, len(urls))])

		if err != nil {
			return tagged, err
		}

		images := jsonRes["outputs"].([]interface{})

		newImages := addImages(images)
		for _, ni := range newImages {
			tagged = append(tagged, ni)
		}
	}
	return tagged, nil
}

func insertRanks(image *TaggedImage, ranks map[string]*ConceptHeap) {
	for name, value := range image.Concept {
		concept := Concept {
			Name: name,
			Image: image,
		}
		
		conceptList, ok := ranks[name]

		if !ok {
			newConcept := &ConceptHeap{ concept }
			heap.Init(newConcept)
			ranks[name] = newConcept
			continue
		}

		if conceptList.Len() < 10 {
			heap.Push(conceptList, concept)
		} else if value > conceptList.Peek() {
			heap.Pop(conceptList)
			heap.Push(conceptList, concept)
		}
	}
}

func (c Client) RankInputs(urls []string) (map[string][]Concept, error) {
	conceptValues := make(map[string][]Concept)
	ranks := make(map[string]*ConceptHeap)
	images, err := c.TagImages(urls)

	if err != nil {
		return conceptValues, err
	}

	for _, img := range images {
		insertRanks(img, ranks)
	}
	
	for concept, conceptHeap := range ranks {
		conceptValues[concept] = Sort(*conceptHeap)
	}

	return conceptValues, nil
}
