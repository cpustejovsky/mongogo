package handlers

import (
	"fmt"
	"net/http"

	"github.com/cpustejovsky/mongogo/helpers"
	"github.com/cpustejovsky/mongogo/internal/models/mongodb/domains"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Logger     *log.Logger
	Collection *mongo.Collection
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *Handler) UpdateDelivered(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get(":domain_name")

	err := domains.UpdateDelivered(h.Collection, domain)
	if err != nil {
		helpers.ServerError(h.Logger, w, err)
		return
	}
	fmt.Fprintf(w, "Successfully updated number of delivered emails for %v", domain)
}

func (h *Handler) UpdateBounced(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get(":domain_name")

	err := domains.UpdateBounced(h.Collection, domain)
	if err != nil {
		helpers.ServerError(h.Logger, w, err)
		return
	}
	fmt.Fprintf(w, "Successfully updated number of bounced emails for %v", domain)
}

func (h *Handler) CheckStatus(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get(":domain_name")

	status, err := domains.CheckStatus(h.Collection, domain)
	if err != nil {
		helpers.ServerError(h.Logger, w, err)
		return
	}
	fmt.Fprintf(w, "Domain %v is status %v", domain, status)
}
