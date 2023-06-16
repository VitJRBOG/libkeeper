package models

// Note stores data about the note.
type Note struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	CreationDate string `json:"c_date"`
}
