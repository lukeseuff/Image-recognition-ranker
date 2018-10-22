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
	taggedImages, concepts, err := apiClient.TagURLs([]string{
		"https://farm5.staticflickr.com/2488/4156338171_54a1df1fde_o.jpg",
		"https://farm8.staticflickr.com/1379/1167031640_c295382ea4_o.jpg",
		"https://c6.staticflickr.com/1/104/303217440_f595f9b310_o.jpg",
		"https://farm5.staticflickr.com/8109/8615989093_dbfd9d5cf5_o.jpg",
		"https://c6.staticflickr.com/4/3532/3811341674_a47fcb3de5_o.jpg",
		"https://farm1.staticflickr.com/3225/2536553703_bf8ec973a2_o.jpg",
		"https://c2.staticflickr.com/6/5565/15116393675_2de9bedcd4_o.jpg",
		"https://c3.staticflickr.com/4/3627/3333969844_7d4dbb971b_o.jpg",
		"https://farm5.staticflickr.com/8203/8153219008_00fbc2842e_o.jpg",
		"https://c6.staticflickr.com/6/5199/7209471856_4a69a064c1_o.jpg",
		"https://c3.staticflickr.com/4/3657/3505128975_5ebd87c04b_o.jpg",
		"https://farm5.staticflickr.com/8228/8451914004_9360928115_o.jpg",
		"https://farm3.staticflickr.com/7105/7143187415_d73938ce15_o.jpg",
		"https://c7.staticflickr.com/9/8399/8691187278_551a7c7819_o.jpg",
		"https://c2.staticflickr.com/3/2851/9227265924_c81f0ba60a_o.jpg",
		"https://c4.staticflickr.com/6/5268/5654895988_a9d875a261_o.jpg",
		"https://farm4.staticflickr.com/4105/5020764789_30bde7ee33_o.jpg",
		"https://farm6.staticflickr.com/3362/3178695285_ce69d74774_o.jpg",
		"https://c1.staticflickr.com/9/8037/8035830054_ec53680ccd_o.jpg",
		"https://c8.staticflickr.com/3/2623/3799825666_335908303e_o.jpg",
		"https://farm5.staticflickr.com/4088/5015123851_1b95877fd0_o.jpg",
	})

	fmt.Printf("%v\n\n", taggedImages)
	for _, concept := range concepts["car"] {
		fmt.Printf("%v - %v\n", concept.Image.Concept["car"], concept.Image.URL)
	}
}
