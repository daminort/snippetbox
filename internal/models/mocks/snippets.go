package mocks

import (
	"time"

	"snippetbox.demien.net/internal/models"
)

var mockedSnippet = models.Snippet{
	ID:      1,
	Title:   "One single snippet",
	Content: "One single snippet\nwalked alone\nin the morning",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetsService struct{}

func (s *SnippetsService) InsertSnippet(title, content string, expires int) (*models.Snippet, error) {
	return &mockedSnippet, nil
}

func (s *SnippetsService) GetSnippet(id int64) (*models.Snippet, error) {
	switch id {
	case 1:
		return &mockedSnippet, nil
	default:
		return nil, models.ErrNotFound
	}
}

func (s *SnippetsService) GetLatest() ([]models.Snippet, error) {
	return []models.Snippet{mockedSnippet}, nil
}
