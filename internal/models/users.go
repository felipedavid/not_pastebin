package models

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID             int64
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPasword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, NOW())`

	_, err = m.DB.Exec(stmt, name, email, hashedPasword)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int64, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int64) (bool, error) {
	return false, nil
}
