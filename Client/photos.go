package client

type Photo struct {
	Id     int32       `json:"id"`
	Width  int32       `json:"width"`
	Height int32       `json:"height"`
	Url    string      `json:"url"`
	Src    PhotoSource `json:"src"`
}
