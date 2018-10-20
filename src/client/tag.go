package client

type ImageData struct {
	Image struct {
		URL string `json:"url"`
	} `json:"image"`
}

type GeneralRequest struct {
	Inputs []struct {
		Data ImageData`json:"data"`
	} `json:"inputs"`
}
