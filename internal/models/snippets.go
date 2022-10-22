package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int64
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m SnippetModel) Get(id int64) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippet WHERE id = $1 AND expires > CURRENT_DATE`

	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return s, nil
}

func (m SnippetModel) Insert(title, content string, expires int64) (int64, error) {
	// TODO: Figure out a way to set a expires statemetn into the SQL query
	stmt := `INSERT INTO snippet (title, content, created, expires) VALUES ($1, $2, NOW(), NOW()) RETURNING id`

	row := m.DB.QueryRow(stmt, title, content)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippet WHERE expires > CURRENT_DATE ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := make([]*Snippet, 0, 10)

	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
