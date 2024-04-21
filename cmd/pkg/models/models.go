package models

import (
	"errors"
)

var ErrNoRecord = errors.New("modes: no matching records found")

type Snippet struct {
	ID	int
	Title string
	Content string
	Created string
	Expires string
}