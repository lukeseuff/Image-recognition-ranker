package main

import (
	// "bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "os"
	"time"
)

type GeneralResponse struct {
	Status struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	} `json:"status"`
	Outputs []struct {
		ID     string `json:"id"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		Model     struct {
			ID         string    `json:"id"`
			Name       string    `json:"name"`
			CreatedAt  time.Time `json:"created_at"`
			AppID      string    `json:"app_id"`
			OutputInfo struct {
				Type    string `json:"type"`
				TypeExt string `json:"type_ext"`
			} `json:"output_info"`
			ModelVersion struct {
				ID        string    `json:"id"`
				CreatedAt time.Time `json:"created_at"`
				Status    struct {
					Code        int    `json:"code"`
					Description string `json:"description"`
				} `json:"status"`
			} `json:"model_version"`
			DisplayName string `json:"display_name"`
		} `json:"model"`
		Input struct {
			ID   string `json:"id"`
			Data struct {
				Image struct {
					URL string `json:"url"`
				} `json:"image"`
			} `json:"data"`
		} `json:"input"`
		Data struct {
			Concepts []struct {
				ID    string  `json:"id"`
				Name  string  `json:"name"`
				Value float64 `json:"value"`
				AppID string  `json:"app_id"`
			} `json:"concepts"`
		} `json:"data"`
	} `json:"outputs"`
}

type Image struct {
	URL string `json:"url"`
}

type Data struct {
	Image `json:"image"`
}

type Inputs struct {
	Data `json:"data"`
}

type GeneralRequest struct {
	Inputs []Inputs `json:"inputs"`
}

func GetPredictions() error {
	var body GeneralRequest
	body.Inputs = append(body.Inputs, Inputs{ Data{ Image: Image{ URL: "https://c5.staticflickr.com/9/8271/8704090235_4da75e857c_o.jpg" } } })
	url := "https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/versions/aa7f35c01e0642fda5cf400f543e7c40/outputs"
	jsonBody, err := json.Marshal(body)

	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))

	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "Key 1694ee43f7754cae9d386ac0eb7a192c")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return err
	}
	
	defer res.Body.Close()
	bdy, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	apiResponse := new(map[string]interface{})
	er := json.Unmarshal(bdy, apiResponse)

	if er != nil {
		return er
	}

	fmt.Printf("%+v\n", apiResponse)
	
	return nil
}

func main() {
	GetPredictions()
	// file, err := os.Open("images.txt")
	// if err != nil {
	// 	return
	// }
	// defer file.Close()
	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }
	// if err := scanner.Err(); err != nil {
	// 	fmt.Printf("%v\n", err)
	// }
}
