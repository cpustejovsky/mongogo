package db

import (
	"time"
)

// Foo represent the structure we need for moving data
// between the app and the database.
type Foo struct {
	ID          string    `db:"user_id"`
	Name        string    `db:"name"`
	Age         string    `db:"email"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}
