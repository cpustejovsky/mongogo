package database

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Set of error variables for CRUD operations.
var (
	ErrDBNotFound = errors.New("not found")
)

type Config struct {
	URI          string
	DatabaseName string
	Ctx          context.Context
}

func Open(cfg Config) (*mongo.Database, error) {
	clientOptions := options.Client().
		ApplyURI(cfg.URI)
	ctx, cancel := context.WithTimeout(cfg.Ctx, 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)
	database := client.Database(cfg.DatabaseName)
	return database, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, client *mongo.Client) error {

	// First check we can ping the database.
	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = client.Ping(ctx, readpref.Primary())
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	// Make sure we didn't timeout or be cancelled.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// TODO: Determine the MongoDB equivalent of this
	// Run a simple query to determine connectivity. Running this query forces a
	// round trip through the database.
	// const q = `SELECT true`
	// var tmp bool
	// return sqlxDB.QueryRowContext(ctx, q).Scan(&tmp)
	return nil
}

// TODO: Determine the MongoDB equivalent of these functions
// // NamedExecContext is a helper function to execute a CUD operation with
// // logging and tracing.
// func NamedExecContext(ctx context.Context, log *zap.SugaredLogger, sqlxDB *sqlx.DB, query string, data interface{}) error {
// 	q := queryString(query, data)
// 	log.Infow("database.NamedExecContext", "traceid", web.GetTraceID(ctx), "query", q)

// 	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "database.query")
// 	span.SetAttributes(attribute.String("query", q))
// 	defer span.End()

// 	if _, err := sqlxDB.NamedExecContext(ctx, query, data); err != nil {
// 		return err
// 	}

// 	return nil
// }

// // NamedQuerySlice is a helper function for executing queries that return a
// // collection of data to be unmarshaled into a slice.
// func NamedQuerySlice(ctx context.Context, log *zap.SugaredLogger, sqlxDB *sqlx.DB, query string, data interface{}, dest interface{}) error {
// 	q := queryString(query, data)
// 	log.Infow("database.NamedQuerySlice", "traceid", web.GetTraceID(ctx), "query", q)

// 	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "database.query")
// 	span.SetAttributes(attribute.String("query", q))
// 	defer span.End()

// 	val := reflect.ValueOf(dest)
// 	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
// 		return errors.New("must provide a pointer to a slice")
// 	}

// 	rows, err := sqlxDB.NamedQueryContext(ctx, query, data)
// 	if err != nil {
// 		return err
// 	}

// 	slice := val.Elem()
// 	for rows.Next() {
// 		v := reflect.New(slice.Type().Elem())
// 		if err := rows.StructScan(v.Interface()); err != nil {
// 			return err
// 		}
// 		slice.Set(reflect.Append(slice, v.Elem()))
// 	}

// 	return nil
// }

// // NamedQueryStruct is a helper function for executing queries that return a
// // single value to be unmarshalled into a struct type.
// func NamedQueryStruct(ctx context.Context, log *zap.SugaredLogger, sqlxDB *sqlx.DB, query string, data interface{}, dest interface{}) error {
// 	q := queryString(query, data)
// 	log.Infow("database.NamedQueryStruct", "traceid", web.GetTraceID(ctx), "query", q)

// 	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "database.query")
// 	span.SetAttributes(attribute.String("query", q))
// 	defer span.End()

// 	rows, err := sqlxDB.NamedQueryContext(ctx, query, data)
// 	if err != nil {
// 		return err
// 	}
// 	if !rows.Next() {
// 		return ErrDBNotFound
// 	}

// 	if err := rows.StructScan(dest); err != nil {
// 		return err
// 	}

// 	return nil
// }

// // queryString provides a pretty print version of the query and parameters.
// func queryString(query string, args ...interface{}) string {
// 	query, params, err := sqlx.Named(query, args)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	for _, param := range params {
// 		var value string
// 		switch v := param.(type) {
// 		case string:
// 			value = fmt.Sprintf("%q", v)
// 		case []byte:
// 			value = fmt.Sprintf("%q", string(v))
// 		default:
// 			value = fmt.Sprintf("%v", v)
// 		}
// 		query = strings.Replace(query, "?", value, 1)
// 	}

// 	query = strings.ReplaceAll(query, "\t", "")
// 	query = strings.ReplaceAll(query, "\n", " ")

// 	return strings.Trim(query, " ")
// }
