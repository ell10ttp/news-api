package models

type Article struct {
	Source      Source `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ImageURL    string `json:"imageUrl"`
	Published   string `json:"published"`
	Content     string `json:"content"`
}
