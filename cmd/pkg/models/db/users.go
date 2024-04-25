package db

import (
	"database/sql"
	"time"

	"github.com/DrEmbryo/snippetbox/cmd/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, ?)`

	createdAt := time.Now().UTC().String();

	result, err := m.DB.Exec(stmt, name, email, string(hashedPassword), createdAt)
	if err != nil{
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)) 
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	
	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}

