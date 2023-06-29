package models

// Version stores data about version of the note.
type Version struct {
	ID           int    `json:"id"`
	FullText     string `json:"full_text"`
	CreationDate string `json:"c_date"`
	Checksum     string `json:"checksum"`
	NoteID       int    `json:"note_id"`
}
