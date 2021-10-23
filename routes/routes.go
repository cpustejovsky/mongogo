package routes

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/cpustejovsky/mongogo/handlers"
	"github.com/cpustejovsky/mongogo/middleware"
	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func Routes(log *log.Logger, client *mongo.Client) http.Handler {

	middlewares := middleware.Middleware{
		Logger: log,
	}

	standardMiddleware := alice.New(middlewares.RecoverPanic, middlewares.SecureHeaders, middlewares.LogRequest)

	mux := pat.New()

	database := client.Database("mongogo")
	collection := database.Collection("users")

	userHandlers := handlers.Handler{
		Logger:     log,
		Collection: collection,
	}
	mux.Get("/api/ping", http.HandlerFunc(userHandlers.Ping))
	mux.Post("/api/user/new", http.HandlerFunc(userHandlers.Create))
	mux.Get("/api/user/:id", http.HandlerFunc(userHandlers.Fetch))
	mux.Put("/api/user/:id", http.HandlerFunc(userHandlers.Update))
	mux.Del("/api/user/:id", http.HandlerFunc(userHandlers.Delete))

	return standardMiddleware.Then(mux)
}
