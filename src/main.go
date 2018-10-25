package main

import (
	"search"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"html/template"
)

type SearchPageData struct {
	Concepts []search.Concept
	Empty    bool
}

func getURLs() []string {
	var urls []string
	file, err := os.Open("data/images.txt")

	if err != nil {
		return urls
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		urls = append(urls, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return urls
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.tmpl"))
	tmpl.Execute(w, "")
}

func main() {
	apiClient, err := search.NewClient()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	concepts, err := apiClient.RankInputs(getURLs()[:2])

	http.HandleFunc("/", handler)
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		c, ok := concepts[query]
		searchPageData := SearchPageData{
			Concepts: c,
			Empty: !ok,
		}
		tmpl := template.Must(template.ParseFiles("template/search.tmpl"))
		tmpl.Execute(w, searchPageData)
	})
	http.ListenAndServe(":8080", nil)
}
