package models

// Version stores data about version of the note.
type Version struct {
	ID           int
	FullText     string
	CreationDate string
	Checksum     string
	NoteID       int
}
