package models

type SourceList struct {
	Sources []Source `json:"sources"`
}
type Source struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"URL"`
	Category    string `json:"category"`
	Language    string `json:"language"`
	Country     string `json:"country"`
}
