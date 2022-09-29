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
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE id = $1 AND expires > CURRENT_DATE`

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

func (m SnippetModel) Insert(title, content string, expires time.Time) (int64, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES ($1, $2, CURRENT_DATE, CURRENT_DATE + INTERVAL "$3 day")`

	row, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets ORDER BY created LIMIT 10`

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
