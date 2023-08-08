package models

// Category stores data about the category.
type Category struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Immutable int    `json:"immutable"`
	IconID    int    `json:"icon_id"`
}
