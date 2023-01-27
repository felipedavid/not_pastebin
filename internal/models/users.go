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
	DB          *sql.DB
	insertStmt  *sql.Stmt
	authStmt    *sql.Stmt
	existsStmt  *sql.Stmt
	getStmt     *sql.Stmt
	getPassStmt *sql.Stmt
	setPassStmt *sql.Stmt
}

func NewUserModel(db *sql.DB) (*UserModel, error) {
	insertStmt, err := db.Prepare(`INSERT INTO users (name, email, hashed_password) VALUES ($1, $2, $3)`)
	if err != nil {
		return nil, err
	}

	authStmt, err := db.Prepare(`SELECT id, hashed_password FROM users WHERE email = $1`)
	if err != nil {
		return nil, err
	}

	existsStmt, err := db.Prepare(`SELECT EXISTS(SELECT true FROM users WHERE id = $1)`)
	if err != nil {
		return nil, err
	}

	getStmt, err := db.Prepare(`SELECT name, email, created from users WHERE id = $1`)
	if err != nil {
		return nil, err
	}

	getPassStmt, err := db.Prepare(`SELECT hashed_password FROM users WHERE id = $1`)
	if err != nil {
		return nil, err
	}

	setPassStmt, err := db.Prepare(`UPDATE users SET hashed_password = $1 WHERE id = $2`)
	if err != nil {
		return nil, err
	}

	return &UserModel{
		DB:          db,
		insertStmt:  insertStmt,
		authStmt:    authStmt,
		existsStmt:  existsStmt,
		getStmt:     getStmt,
		getPassStmt: getPassStmt,
		setPassStmt: setPassStmt,
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

	err := m.authStmt.QueryRow(email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidEmail
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidPassword
		}
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	err := m.existsStmt.QueryRow(id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	var user User

	err := m.getStmt.QueryRow(id).Scan(&user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) PasswordMatch(id int, password string) (bool, error) {
	var hashedPassword []byte

	err := m.getPassStmt.QueryRow(id).Scan(&hashedPassword)
	if err != nil {
		return false, nil
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, ErrInvalidPassword
		}
		return false, err
	}

	return true, nil
}

func (m *UserModel) SetPassword(id int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	_, err = m.setPassStmt.Exec(id, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
