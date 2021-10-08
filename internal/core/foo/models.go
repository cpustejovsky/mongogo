package foo

import (
	"time"
	"unsafe"

	"github.com/ardanlabs/service/business/core/user/db"
)

// Foo represents an individual user.
type Foo struct {
	Name        string    `json:"id"`
	Age         int       `json:"age"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// NewFoo contains information needed to create a new Foo.
type NewFoo struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required,age"`
}

// UpdateFoo defines what information may be provided to modify an existing
// Foo. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateFoo struct {
	Name *string `json:"name"`
	Age  *int    `json:"age" validate:"omitempty,age"`
}

// =============================================================================

func toFoo(dbUsr db.Foo) Foo {
	pu := (*Foo)(unsafe.Pointer(&dbUsr))
	return *pu
}

func toFooSlice(dbUsrs []db.Foo) []Foo {
	users := make([]Foo, len(dbUsrs))
	for i, dbUsr := range dbUsrs {
		users[i] = toFoo(dbUsr)
	}
	return users
}
