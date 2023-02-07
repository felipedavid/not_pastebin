package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB         *sql.DB
	insertStmt *sql.Stmt
	getStmt    *sql.Stmt
	latestStmt *sql.Stmt
}

// NewSnippetModel allocates a new SnippetModel object and initializes it
// with all pre-compiled statements that query the snippets table
func NewSnippetModel(db *sql.DB) (*SnippetModel, error) {
	// TODO: Fix expiration date
	insertStmt, err := db.Prepare(`INSERT INTO snippets (title, content, expires) VALUES
		($1, $2, $3) RETURNING id`)
	if err != nil {
		return nil, err
	}

	getStmt, err := db.Prepare(`SELECT id, title, content, created, expires FROM snippets 
		WHERE id = $1 AND expires > NOW()`)
	if err != nil {
		return nil, err
	}

	latestStmt, err := db.Prepare(`SELECT id, title, content, created, expires FROM snippets
		WHERE expires > NOW() ORDER BY id DESC LIMIT 10`)
	if err != nil {
		return nil, err
	}

	return &SnippetModel{
		DB:         db,
		insertStmt: insertStmt,
		getStmt:    getStmt,
		latestStmt: latestStmt,
	}, nil
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	id := 0
	err := m.insertStmt.QueryRow(title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	s := &Snippet{}
	err := m.getStmt.QueryRow(id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Why we don't return sql.ErrNoRows directly? We want to isolate our model
			// from the database details, so we can change our database system without
			// hassle
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	rows, err := m.latestStmt.Query()
	if err != nil {
		return nil, err
	}
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
