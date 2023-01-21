package models

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB         *sql.DB
	insertStmt *sql.Stmt
}

func NewUserModel(db *sql.DB) (*UserModel, error) {
	insertStmt, err := db.Prepare(`INSERT INTO users (name, email, hashed_password) VALUES ($1, $2, $3)`)
	if err != nil {
		return nil, err
	}

	return &UserModel{
		DB:         db,
		insertStmt: insertStmt,
	}, nil
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = m.insertStmt.Exec(name, email, hashedPassword)
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23505" && strings.Contains(pqError.Message, "users_email_uc") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
