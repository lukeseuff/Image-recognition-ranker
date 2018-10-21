package main

import (
	"client"
	"bufio"
	// "bytes"
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	// "net/http"
	"os"
	// "time"
)

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

func main() {
	apiClient, err := client.NewClient()
	
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	// apiClient.TagURL("https://c3.staticflickr.com/6/5615/15629176482_0384ab8a9f_o.jpg")
	// taggedImages, err := apiClient.TagURLs(getURLs()[:10])
	taggedImages, concepts, err := apiClient.TagURLs([]string{"https://farm5.staticflickr.com/2488/4156338171_54a1df1fde_o.jpg", "https://farm8.staticflickr.com/1379/1167031640_c295382ea4_o.jpg", "https://c6.staticflickr.com/1/104/303217440_f595f9b310_o.jpg"})
	fmt.Printf("%v\n\n", taggedImages)
	fmt.Printf("%v\n", concepts)
}
