package server

import (
	"request"
	"rank"
	"util"
	"net/http"
	"tag"
	"html/template"
)

const maxBatchSize = 128

type SearchPageData struct {
	Concepts []tag.Concept
	Empty    bool
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.tmpl"))
	tmpl.Execute(w, "")
}

func search(concepts map[string][]tag.Concept) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		c, ok := concepts[query]
		searchPageData := SearchPageData{
			Concepts: c,
			Empty: !ok,
		}
		tmpl := template.Must(template.ParseFiles("template/search.tmpl"))
		tmpl.Execute(w, searchPageData)
	}
}

func Start() {
	urls := util.GetURLs("data/images.txt")
	apiClient, _ := request.NewClient()
	responses := apiClient.BatchPrediction(urls, maxBatchSize)
	taggedImages := tag.TagImages(responses)
	rankings := rank.RankTaggedImages(taggedImages)

	http.HandleFunc("/", index)
	http.HandleFunc("/search", search(rankings))
	http.ListenAndServe(":8080", nil)
}
