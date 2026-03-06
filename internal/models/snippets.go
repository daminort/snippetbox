package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetServiceInstance interface {
	InsertSnippet(title, content string, expires int) (*Snippet, error)
	GetSnippet(id int64) (*Snippet, error)
	GetLatest() ([]Snippet, error)
}

type Snippet struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

type SnippetsService struct {
	DB *sql.DB
}

func NewSnippetsService(db *sql.DB) *SnippetsService {
	return &SnippetsService{
		DB: db,
	}
}

func (s *SnippetsService) InsertSnippet(title, content string, expires int) (*Snippet, error) {
	query := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
	`

	result, err := s.DB.Exec(query, title, content, expires)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	snippet := Snippet{
		ID:      id,
		Title:   title,
		Content: content,
		Created: time.Now(),
		Expires: time.Now().AddDate(0, 0, expires),
	}

	return &snippet, nil
}

func (s *SnippetsService) GetSnippet(id int64) (*Snippet, error) {
	query := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND id = ?
	`

	var sn Snippet

	row := s.DB.QueryRow(query, id)
	err := row.Scan(&sn.ID, &sn.Title, &sn.Content, &sn.Created, &sn.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &sn, nil
}

func (s *SnippetsService) GetLatest() ([]Snippet, error) {
	query := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP()
		ORDER BY id DESC
		LIMIT 10
	`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var sn Snippet
		err := rows.Scan(&sn.ID, &sn.Title, &sn.Content, &sn.Created, &sn.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, sn)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
