package client

type Data struct {}

type VideoData struct {
	Data
	Video struct {
		URL string `json:"url"`
	} `json:"video"`
}

type ImageData struct {
	Data
	Image struct {
		URL string `json:"url"`
	} `json:"image"`
}

type GeneralRequest struct {
	Inputs []struct {
		Data Data `json:"data"`
	} `json:"inputs"`
}
