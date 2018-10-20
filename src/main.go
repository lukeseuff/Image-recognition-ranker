package main

import (
	"client"
	// "bufio"
	// "bytes"
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	// "net/http"
	"os"
	// "time"
)

func main() {
	apiClient, err := client.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	apiClient.TagUrl("https://c3.staticflickr.com/6/5615/15629176482_0384ab8a9f_o.jpg")
}
