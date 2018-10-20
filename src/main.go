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

func getUrls() []string {
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
	
	// apiClient.TagUrl("https://c3.staticflickr.com/6/5615/15629176482_0384ab8a9f_o.jpg")

	apiClient.TagUrls(getUrls()[:2])
}
