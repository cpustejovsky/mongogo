package handlers

import (
	"net/http"

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	//get JSON body and decode
	//create new document within mongodb table
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	//get id from url
	//find user by id and return
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	//get id from url
	//find and update user with id
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	//get id from url
	//find and delete user with id
}
