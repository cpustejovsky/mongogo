// Package db contains foo related CRUD functionality.
package db

import (
	"context"
	"fmt"

	"github.com/cpustejovsky/mongogo/internal/sys/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Store manages the set of API's for foo access.
type Store struct {
	log     *zap.SugaredLogger
	mongoDB *mongo.Database
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, mongoDB *mongo.Database) Store {
	return Store{
		log:     log,
		mongoDB: mongoDB,
	}
}

// Create inserts a new foo into the database.
func (s Store) Create(ctx context.Context, foo Foo) error {
	const q = `
	INSERT INTO foos
		(foo_id, name, email, password_hash, roles, date_created, date_updated)
	VALUES
		(:foo_id, :name, :email, :password_hash, :roles, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.mongoDB, q, foo); err != nil {
		return fmt.Errorf("inserting foo: %w", err)
	}


	return nil
}

// Update replaces a foo document in the database.
func (s Store) Update(ctx context.Context, foo Foo) error {

	s.mongoDB.Collection()

	const q = `
	UPDATE
		foos
	SET 
		"name" = :name,
		"email" = :email,
		"roles" = :roles,
		"password_hash" = :password_hash,
		"date_updated" = :date_updated
	WHERE
		foo_id = :foo_id`

	if err := database.NamedExecContext(ctx, s.log, s.mongoDB, q, foo); err != nil {
		return fmt.Errorf("updating fooID[%s]: %w", foo.ID, err)
	}

	return nil
}

// Delete removes a foo from the database.
func (s Store) Delete(ctx context.Context, fooID string) error {
	data := struct {
		FooID string `db:"foo_id"`
	}{
		FooID: fooID,
	}

	const q = `
	DELETE FROM
		foos
	WHERE
		foo_id = :foo_id`

	if err := database.NamedExecContext(ctx, s.log, s.mongoDB, q, data); err != nil {
		return fmt.Errorf("deleting fooID[%s]: %w", fooID, err)
	}

	return nil
}

// Query retrieves a list of existing foos from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Foo, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		foos
	ORDER BY
		foo_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var foos []Foo
	if err := database.NamedQuerySlice(ctx, s.log, s.mongoDB, q, data, &foos); err != nil {
		return nil, fmt.Errorf("selecting foos: %w", err)
	}

	return foos, nil
}

// QueryByID gets the specified foo from the database.
func (s Store) QueryByID(ctx context.Context, fooID string) (Foo, error) {
	data := struct {
		FooID string `db:"foo_id"`
	}{
		FooID: fooID,
	}

	const q = `
	SELECT
		*
	FROM
		foos
	WHERE 
		foo_id = :foo_id`

	var foo Foo
	if err := database.NamedQueryStruct(ctx, s.log, s.mongoDB, q, data, &foo); err != nil {
		return Foo{}, fmt.Errorf("selecting fooID[%q]: %w", fooID, err)
	}

	return foo, nil
}

// QueryByEmail gets the specified foo from the database by email.
func (s Store) QueryByEmail(ctx context.Context, email string) (Foo, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	const q = `
	SELECT
		*
	FROM
		foos
	WHERE
		email = :email`

	var foo Foo
	if err := database.NamedQueryStruct(ctx, s.log, s.mongoDB, q, data, &foo); err != nil {
		return Foo{}, fmt.Errorf("selecting email[%q]: %w", email, err)
	}

	return foo, nil
}
