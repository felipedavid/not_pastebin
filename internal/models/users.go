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
	getStmt    *sql.Stmt
}

func NewUserModel(db *sql.DB) (*UserModel, error) {
	insertStmt, err := db.Prepare(`INSERT INTO users (name, email, hashed_password) VALUES ($1, $2, $3)`)
	if err != nil {
		return nil, err
	}

	getStmt, err := db.Prepare(`SELECT id, hashed_password FROM users WHERE email = $1`)
	if err != nil {
		return nil, err
	}

	return &UserModel{
		DB:         db,
		insertStmt: insertStmt,
		getStmt:    getStmt,
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
	var id int
	var hashedPassword []byte

	err := m.getStmt.QueryRow(email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
