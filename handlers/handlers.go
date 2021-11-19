package handlers

import (
	"fmt"
	"net/http"

	"github.com/cpustejovsky/mongogo/helpers"
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

func (h *Handler) PingPanic(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("requestId")
	idstr := fmt.Sprintf("Request ID: %v\n", id)
	w.Write([]byte(idstr))
	panic("foo")
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	//get JSON body and decode
	fUser, err := helpers.DecodeUserForm(r)
	if err != nil {
		if err == helpers.EmptyBodyError {
			fmt.Fprint(w, err)
			return
		}
		helpers.ServerError(h.Logger, w, err)
		return
	}
	missingProperties := []string{}
	if fUser.Name == nil {
		missingProperties = append(missingProperties, "Name")
	}
	if fUser.Age == nil {
		missingProperties = append(missingProperties, "Age")
	}
	if fUser.Email == nil {
		missingProperties = append(missingProperties, "Email")
	}
	if len(missingProperties) > 0 {
		fmt.Fprint(w, helpers.MissingPropertyError(missingProperties))
		return
	}

	fmt.Fprint(w, *fUser.Name)
	//create new document within mongodb table
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	//get id from url
	//find user by id and return
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	//get id from url
	//get JSON body and decode
	//find and update user with id
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	//get id from url
	//find and delete user with id
}
