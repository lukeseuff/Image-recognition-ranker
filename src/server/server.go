package server

import (
	"request"
	"rank"
	"util"
	"net/http"
	"runtime"
	"os/exec"
	"tag"
	"html/template"
	"fmt"
)

const maxBatchSize = 128

type SearchPageData struct {
	Concepts []tag.Concept
	Empty    bool
}

func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
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

	startBrowser("http://localhost:8080")
	fmt.Println("serving on localhost:8080")
	http.HandleFunc("/", index)
	http.HandleFunc("/search", search(rankings))
	http.ListenAndServe(":8080", nil)
}
