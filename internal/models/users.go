package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UsersServiceInstance interface {
	GetUser(id int64) (*User, error)
	InsertUser(name, email, password string) (*User, error)
	Authenticate(email, password string) (*User, error)
}

type User struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
}

type UsersService struct {
	DB *sql.DB
}

func NewUsersService(db *sql.DB) *UsersService {
	return &UsersService{
		DB: db,
	}
}

func (s *UsersService) GetUser(id int64) (*User, error) {
	query := `
		SELECT id, name, email, created
		FROM users
		WHERE id = ?
	`

	var user User

	row := s.DB.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *UsersService) InsertUser(name, email, password string) (*User, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (name, email, password, created)
		VALUES (?, ?, ?, UTC_TIMESTAMP())
	`

	result, err := s.DB.Exec(query, name, email, string(hashedPwd))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return nil, ErrDuplicateEmail
			}
		}

		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetUser(userID)
}

func (s *UsersService) Authenticate(email, password string) (*User, error) {
	query := `
		SELECT id, name, email, password, created
		FROM users
		WHERE email = ?
	`

	var u User
	var pwd string

	row := s.DB.QueryRow(query, email)
	err := row.Scan(&u.ID, &u.Name, &u.Email, &pwd, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	hash := []byte(pwd)
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return &u, nil
}
