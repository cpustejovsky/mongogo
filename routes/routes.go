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

	database := client.Database("mongogo_domains")
	collection := database.Collection("domains")

	domainHandlers := handlers.Handler{
		Logger:     log,
		Collection: collection,
	}
	mux.Put("/events/:domain_name/delivered", standardMiddleware.ThenFunc(domainHandlers.UpdateDelivered))
	mux.Put("/events/:domain_name/bounced", standardMiddleware.ThenFunc(domainHandlers.UpdateBounced))
	mux.Get("/domains/:domain_name", standardMiddleware.ThenFunc(domainHandlers.CheckStatus))
	mux.Get("/ping", standardMiddleware.ThenFunc(domainHandlers.Ping))

	return mux
}
