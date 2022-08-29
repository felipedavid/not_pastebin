package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UsersModel struct {
	DB *sql.DB
}

func (m *UsersModel) Insert(name, email, password string) error {
	return nil
}

func (m *UsersModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UsersModel) Exists(id int) (bool, error) {
	return false, nil
}
