package main

import (
	"client"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"html/template"
)

type SearchPageData struct {
	Concepts []client.Concept
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
	apiClient, err := client.NewClient()
	
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	taggedImages, concepts, err := apiClient.TagURLs(getURLs()[0:256])

	fmt.Printf("%v\n\n", taggedImages)

	http.HandleFunc("/", handler)
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		searchPageData := SearchPageData{ Concepts: concepts[query] }
		tmpl := template.Must(template.ParseFiles("template/search.tmpl"))
		tmpl.Execute(w, searchPageData)
	})
	http.ListenAndServe(":8080", nil)
}
