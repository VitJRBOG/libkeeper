package models

// Icon stores data about the icons.
type Icon struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Path string `json:"path"`
}
