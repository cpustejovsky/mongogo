package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
)

type Domain struct {
	Name      string ``
	Bounced   int    ``
	Delivered int    ``
}
