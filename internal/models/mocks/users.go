package mocks

import (
	"time"

	"snippetbox.demien.net/internal/models"
)

var mockedUser = models.User{
	ID:      1,
	Name:    "Sam",
	Email:   "sam@gmail.com",
	Created: time.Now(),
}

type UsersService struct{}

func (s *UsersService) GetUser(id int64) (*models.User, error) {
	switch id {
	case 1:
		return &mockedUser, nil
	default:
		return nil, models.ErrNotFound
	}
}

func (s *UsersService) InsertUser(name, email, password string) (*models.User, error) {
	switch email {
	case "admin@gmail.com":
		return nil, models.ErrDuplicateEmail
	default:
		return &mockedUser, nil
	}
}

func (s *UsersService) Authenticate(email, password string) (*models.User, error) {
	switch email {
	case "sam@gmail.com":
		return &mockedUser, nil
	default:
		return nil, models.ErrInvalidCredentials
	}
}
