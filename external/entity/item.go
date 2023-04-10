package entity

// Item represents a business returned by the events events searching API.
type Item struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Rating     float64  `json:"rating"`
	Address    string   `json:"address"`
	Categories []string `json:"categories"`
	ImageURL   string   `json:"image_url"`
	Url        string   `json:"url"`
	Distance   float64  `json:"distance"`
}
