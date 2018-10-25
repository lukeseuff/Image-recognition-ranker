package rank

import (
	"container/heap"
	"tag"
)

type ConceptHeap []tag.Concept

func (h ConceptHeap) Len()          int     { return len(h) }
func (h ConceptHeap) Less(i, j int) bool    { return h[i].GetValue() < h[j].GetValue() }
func (h ConceptHeap) Swap(i, j int)         { h[i], h[j] = h[j], h[i] }
func (h ConceptHeap) Peek()         float64 { return h[0].GetValue() }

func Sort(h ConceptHeap) []tag.Concept {
	sorted := make([]tag.Concept, len(h))
	length := len(h)
	for i := 0; i < length; i++ {
		sorted[length - i - 1] = heap.Pop(&h).(tag.Concept)
	}
	return sorted
}

func (h *ConceptHeap) Push(x interface{}) {
	*h = append(*h, x.(tag.Concept))
}

func (h *ConceptHeap) Pop() interface {} {
	old := *h
	n := len(old)
	x := old[n - 1]
	*h = old[0 : n - 1]
	return x
}

func insertRanks(image *tag.TaggedImage, ranks map[string]*ConceptHeap) {
	for name, value := range image.Concept {
		concept := tag.Concept {
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


func RankTaggedImages(images []*tag.TaggedImage) map[string][]tag.Concept {
	conceptValues := make(map[string][]tag.Concept)
	ranks := make(map[string]*ConceptHeap)

	for _, img := range images {
		insertRanks(img, ranks)
	}
	
	for concept, conceptHeap := range ranks {
		conceptValues[concept] = Sort(*conceptHeap)
	}

	return conceptValues
}
