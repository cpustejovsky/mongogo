package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
)

type User struct {
	Name   string ``
	Email  string ``
	Age    int    ``
	Active bool   ``
}
