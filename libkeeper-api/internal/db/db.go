package db

import (
	"database/sql"
	"fmt"
	"libkeeper-api/internal/models"

	_ "github.com/lib/pq" // Postgres driver
)

// Connection stores DB connection.
type Connection struct {
	Conn *sql.DB
}

// NewConnection creates a connection to the PostgreSQL database and returns the struct with it.
func NewConnection(dsn string) (Connection, error) {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return Connection{}, fmt.Errorf("unable to create a database connection: %s", err.Error())
	}

	err = dbConn.Ping()
	if err != nil {
		return Connection{}, fmt.Errorf("unable connect to database: %s", err.Error())
	}

	return Connection{
		Conn: dbConn,
	}, nil
}

// SelectIcons selects entries from the "icon" table.
func SelectIcons(dbConn Connection) ([]models.Icon, error) {
	query := "SELECT * FROM icon"

	rows, err := dbConn.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to select entries from the 'icon' table: %s", err)
	}

	icons := []models.Icon{}

	for rows.Next() {
		icon := models.Icon{}

		err := rows.Scan(&icon.ID, &icon.Type, &icon.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows from the 'icon' table: %s", err)
		}

		icons = append(icons, icon)
	}

	return icons, nil
}

// CreateCategory inserts new entry into the "category" table.
func CreateCategory(dbConn Connection, category models.Category) error {
	query := "INSERT INTO category(name, icon_id) VALUES($1, $2)"

	_, err := dbConn.Conn.Exec(query, category.Name, category.IconID)
	if err != nil {
		return fmt.Errorf("failed to insert entries into the 'category' table: %s", err)
	}

	return nil
}

// SelectCategories selects entries from the "category" table.
func SelectCategories(dbConn Connection) ([]models.Category, error) {
	query := "SELECT * FROM category"

	rows, err := dbConn.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to select entries from the 'category' table: %s", err)
	}

	categories := []models.Category{}

	for rows.Next() {
		category := models.Category{}

		err := rows.Scan(&category.ID, &category.Name, &category.Immutable, &category.IconID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows from the 'category' table: %s", err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// UpdateCategory updates an existing entry in the 'category' table.
func UpdateCategory(dbConn Connection, category models.Category) error {
	query := "UPDATE category SET name = $1, icon_id = $2 WHERE id = $3 AND immutable = '0'"

	_, err := dbConn.Conn.Exec(query, category.Name, category.IconID, category.ID)
	if err != nil {
		return fmt.Errorf("failed to update the 'category' table entry: %s", err)
	}

	return nil
}

// DeleteCategory deletes an existing entry from the 'category' table if its 'removable' field is '1'.
func DeleteCategory(dbConn Connection, categoryID int) error {
	query := "DELETE FROM category WHERE id = $1 AND immutable = '0'"

	_, err := dbConn.Conn.Exec(query, categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete the 'category' table entry: %s", err)
	}

	return nil
}

// CreateNote inserts new entries into the "note" and "version" tables.
func CreateNote(dbConn Connection, note models.Note, version models.Version) error {
	beginOfQuery := "WITH new_note AS (INSERT INTO note(title, c_date) VALUES($1, $2) RETURNING id), " +
		"new_note_category AS (INSERT INTO note_category(note_id, category_id) " +
		"SELECT * FROM ("

	middlePartOfQuery := ""
	values := []interface{}{
		note.Title, note.CreationDate,
	}
	n := 3

	for i, category := range note.Categories {
		a := fmt.Sprintf("SELECT (SELECT id FROM new_note) AS note_id, (SELECT id FROM category WHERE name=$%d) AS category_id",
			n)
		n++
		middlePartOfQuery += a
		if i < len(note.Categories)-1 {
			middlePartOfQuery += " UNION ALL "
		}
		values = append(values, category)
	}

	endOfQuery := fmt.Sprintf(") AS new_note_categories) "+
		"INSERT INTO version(full_text, c_date, checksum, note_id) "+
		"VALUES($%d, $%d, $%d, (SELECT id FROM new_note))", n, n+1, n+2)

	values = append(values, version.FullText, version.CreationDate, version.Checksum)

	query := beginOfQuery + middlePartOfQuery + endOfQuery

	stmt, err := dbConn.Conn.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare the query: %s", err)
	}

	_, err = stmt.Exec(values...)
	if err != nil {
		return fmt.Errorf("failed to insert entries into the 'note', 'version' and 'note_category' tables: %s", err)
	}

	return nil
}

// SelectNotes selects entries from the "note" table.
func SelectNotes(dbConn Connection) ([]models.Note, error) {
	query := "SELECT note.id, note.title, note.c_date, category.name FROM note " +
		"INNER JOIN note_category ON note.id = note_category.note_id " +
		"INNER JOIN category ON category.id = note_category.category_id " +
		"ORDER BY note.id ASC"

	rows, err := dbConn.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to select entries from the 'note' table: %s", err)
	}

	notes := []models.Note{}
	lastID := 0

	for rows.Next() {
		note := models.Note{}
		category := ""

		err := rows.Scan(&note.ID, &note.Title, &note.CreationDate, &category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows from the 'note' table: %s", err)
		}

		if note.ID == lastID {
			notes[len(notes)-1].Categories = append(notes[len(notes)-1].Categories, category)
		} else {
			note.Categories = append(note.Categories, category)
			notes = append(notes, note)
		}

		lastID = note.ID
	}

	return notes, nil
}

// UpdateNoteCategories inserts new records and deletes old ones by 'note_id' in the 'note_category' table.
func UpdateNoteCategories(dbConn Connection, note models.Note) error {
	beginOfQuery := "WITH new_relations AS ("
	endOfQuery := "), create_relations AS (INSERT INTO note_category (note_id, category_id) " +
		"SELECT * FROM new_relations WHERE NOT EXISTS (" +
		"SELECT id FROM note_category " +
		"WHERE note_category.note_id = new_relations.note_id AND " +
		"note_category.category_id = new_relations.category_id)) " +
		"DELETE FROM note_category " +
		"WHERE note_id IN (SELECT note_id FROM new_relations) AND category_id NOT IN " +
		"(SELECT category_id FROM new_relations)"

	values := []interface{}{}
	queryWithValues := ""
	n := 1

	for i, category := range note.Categories {
		queryWithValues += fmt.Sprintf("SELECT (SELECT id FROM note WHERE id = $%d) AS note_id, "+
			"(SELECT id FROM category WHERE name = $%d) AS category_id",
			n, n+1)
		if i < len(note.Categories)-1 {
			queryWithValues += " UNION ALL "
		}

		n += 2
		values = append(values, note.ID, category)
	}

	query := beginOfQuery + queryWithValues + endOfQuery

	_, err := dbConn.Conn.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to update the 'note_category' entries: %s", err)
	}

	return nil
}

// SelectVersions selects entries from the "version" table by the value of "note_id" field. Returns them sorted by DESC.
func SelectVersions(dbConn Connection, noteID int) ([]models.Version, error) {
	query := "SELECT * FROM version WHERE note_id = $1 ORDER BY c_date DESC"

	rows, err := dbConn.Conn.Query(query, noteID)
	if err != nil {
		return nil, fmt.Errorf("failed to select entries from the 'version' table: %s", err)
	}

	versions := []models.Version{}

	for rows.Next() {
		version := models.Version{}

		err := rows.Scan(&version.ID, &version.FullText, &version.CreationDate, &version.Checksum, &version.NoteID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows from the 'version' table: %s", err)
		}

		versions = append(versions, version)
	}

	return versions, nil
}

// UpdateNote updates an existing entry in the 'note' table and inserts a new entry in the 'version' table.
func UpdateNote(dbConn Connection, note models.Note, version models.Version) error {
	transaction, err := dbConn.Conn.Begin()

	defer transaction.Rollback()

	result, err := transaction.Exec("UPDATE note SET title=$1 WHERE id=$2", note.Title, note.ID)
	if err != nil {
		return fmt.Errorf("failed to update the 'note' table entry: %s", err.Error())
	}

	rowsNumber, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %s", err.Error())
	}
	if rowsNumber == 0 {
		return fmt.Errorf("no notes found for the %d ID", note.ID)
	}

	rows, err := transaction.Query("SELECT id FROM version WHERE note_id = $1::integer AND checksum = $2",
		note.ID, version.Checksum)

	versionID := -1

	for rows.Next() {
		err := rows.Scan(&versionID)
		if err != nil {
			return fmt.Errorf("failed to scan rows from the 'version' table: %s", err)
		}
	}

	if versionID == -1 {
		query := "INSERT INTO version(full_text, c_date, checksum, note_id) VALUES($1, $2, $3, $4::integer)"

		result, err = transaction.Exec(query,
			version.FullText, version.CreationDate, version.Checksum, note.ID)
		if err != nil {
			return fmt.Errorf("failed to create the new 'version' table entry: %s", err.Error())
		}

		rowsNumber, err = result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get affected rows: %s", err.Error())
		}
		if rowsNumber == 0 {
			return fmt.Errorf("no versions created for the %d note ID", note.ID)
		}
	} else {
		query := "UPDATE version SET c_date = $1 WHERE id = $2::integer"
		result, err := transaction.Exec(query, version.CreationDate, versionID)
		if err != nil {
			return fmt.Errorf("failed to update the 'version' table entry: %s", err.Error())
		}

		rowsNumber, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get affected rows: %s", err.Error())
		}
		if rowsNumber == 0 {
			return fmt.Errorf("no versions found for the %d ID", versionID)
		}
	}

	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit the transaction: %s", err.Error())
	}

	return nil
}

// DeleteNote deletes an existing entry from the 'note' table and deletes the associated entries from the 'version' table.
func DeleteNote(dbConn Connection, noteID int) error {
	query := "WITH deleted_versions AS (DELETE FROM version WHERE note_id = $1 RETURNING note_id) " +
		"DELETE FROM note WHERE id IN (SELECT note_id FROM deleted_versions)"

	_, err := dbConn.Conn.Exec(query, noteID)
	if err != nil {
		return fmt.Errorf("failed to delete the 'note' table entry: %s", err)
	}

	return nil
}
