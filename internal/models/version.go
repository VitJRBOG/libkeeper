package models

// Version store a data from 'versions' table.
type Version struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Date     int64  `json:"date"`
	Checksum string `json:"checksum"`
	NoteID   int    `json:"note_id"`
}
