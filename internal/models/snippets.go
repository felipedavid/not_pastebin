package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet represents a row in the snippets database table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel is just a wrapper around a database connection pool with methods to manipulate
// the snippets table
type SnippetModel struct {
	DB *sql.DB
}

// Insert inserts a snippet into the database and returns its id
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	r, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get returns a snippet with the specified id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
        WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query has not rows, then row.Scan() will return a
		// sql.ErrNoRows error. We check for the error and then
		// return our own. Why we are doing this? To encapsulate
		// the model, and make it easy to change the data storage
		// without changing the models interface later on
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

// Latest returns the 10 last inserted snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
        WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Closes the result set. While the resultset is not closed, the connection
	// to the database will remain open
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	return snippets, nil
}
