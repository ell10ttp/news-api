package models

type SourceList struct {
	Sources []Source `json:"sources"`
}
type Source struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	URL         string            `json:"URL"`
	Language    string            `json:"language"`
	Country     string            `json:"country"`
	Categories  map[string]string `json:"categories"`
}
