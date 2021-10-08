// Package foo provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package foo

import (
	"context"
	"fmt"
	"time"

	"github.com/cpustejovsky/monggog/internal/core/foo/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Core manages the set of API's for foo access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for foo api access.
func NewCore(log *zap.SugaredLogger, mongoDB *mongo.Database) Core {
	return Core{
		store: db.NewStore(log, mongoDB),
	}
}

// Create inserts a new foo into the database.
func (c Core) Create(ctx context.Context, nu NewFoo, now time.Time) (Foo, error) {

	dbFoo := db.Foo{
		Name: nu.Name,
		Age:  nu.Age,
	}

	if err := c.store.Create(ctx, dbFoo); err != nil {
		return Foo{}, fmt.Errorf("create: %w", err)
	}

	return toFoo(dbFoo), nil
}

// Update replaces a foo document in the database.
func (c Core) Update(ctx context.Context, fooID string, uu UpdateFoo, now time.Time) error {

	dbFoo, err := c.store.QueryByID(ctx, fooID)
	if err != nil {
		return fmt.Errorf("updating foo fooID[%s]: %w", fooID, err)
	}

	if uu.Name != nil {
		dbFoo.Name = *uu.Name
	}
	dbFoo.DateUpdated = now

	if err := c.store.Update(ctx, dbFoo); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// Delete removes a foo from the database.
func (c Core) Delete(ctx context.Context, fooID string) error {
	if err := c.store.Delete(ctx, fooID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing foos from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Foo, error) {
	dbFoos, err := c.store.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toFooSlice(dbFoos), nil
}

// QueryByID gets the specified foo from the database.
func (c Core) QueryByID(ctx context.Context, fooID string) (Foo, error) {

	dbFoo, err := c.store.QueryByID(ctx, fooID)
	if err != nil {
		return Foo{}, fmt.Errorf("query: %w", err)
	}

	return toFoo(dbFoo), nil
}

// QueryByEmail gets the specified foo from the database by email.
func (c Core) QueryByEmail(ctx context.Context, email string) (Foo, error) {

	dbFoo, err := c.store.QueryByEmail(ctx, email)
	if err != nil {
		return Foo{}, fmt.Errorf("query: %w", err)
	}

	return toFoo(dbFoo), nil
}
