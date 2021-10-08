// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"net/http"

	"github.com/ardanlabs/service/foundation/web"
	"github.com/cpustejovsky/mongogo/app/default/handlers/foogrp"
	"github.com/cpustejovsky/mongogo/internal/core/foo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log *zap.SugaredLogger
	DB  *mongo.Database
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	// Register user management and authentication endpoints.
	fgh := foogrp.Handlers{
		Core: foo.NewCore(cfg.Log, cfg.DB),
	}
	app.Handle(http.MethodPost, version, "/api/foogrp/:id", fgh.Create)
	app.Handle(http.MethodGet, version, "/api/foogrp/", fgh.Query)
	app.Handle(http.MethodGet, version, "/api/foogrp/:id", fgh.QueryByID)
	app.Handle(http.MethodPut, version, "/api/foogrp/:id", fgh.Update)
	app.Handle(http.MethodDelete, version, "/api/foogrp/:id", fgh.Delete)

}
