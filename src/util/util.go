package util

import (
	"bufio"
	"os"
  "encoding/json"
  "io/ioutil"
)

func GetConfig() map[string]interface{} {
  var config map[string]interface{}
  configFile, _ := ioutil.ReadFile("config/config.json")
  json.Unmarshal(configFile, &config)

  return config
}

func GetURLs(filename string) []string {
	var urls []string
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		urls = append(urls, line)
	}

  file.Close()

	return urls
}
