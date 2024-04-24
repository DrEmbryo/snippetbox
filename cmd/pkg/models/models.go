package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("modes: no matching records found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail = errors.New("models: duplicate email")
) 

type Snippet struct {
	ID	int
	Title string
	Content string
	Created string
	Expires string
}

type User struct {
	ID int
	Name string
	Email string
	HashedPassword string
	Created string
}