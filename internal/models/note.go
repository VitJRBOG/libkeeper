package models

// Note store a data from 'notes' table.
type Note struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Date  int64  `json:"date"`
}
