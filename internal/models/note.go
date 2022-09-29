package models

import "time"

// Note store a data from 'notes' table.
type Note struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
}
