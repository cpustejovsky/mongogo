package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cpustejovsky/mongogo/helpers"
	"github.com/cpustejovsky/mongogo/internal/models/mongodb/user"
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
	//create new document within mongodb table
	user, err := user.Create(h.Collection, fUser)
	if err != nil {
		fmt.Fprint(w, errors.New("Unable to Update Item"))
		return
	}
	fmt.Fprint(w, user)
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	//get id from url
	id := strings.TrimPrefix(r.URL.Path, "/api/user/")
	//find user by id and return
	user, err := user.Fetch(h.Collection, id)
	if err != nil {
		fmt.Fprint(w, errors.New("Unable to Fetch Item"))
		return
	}
	fmt.Fprint(w, user)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	//get id from url
	id := strings.TrimPrefix(r.URL.Path, "/api/user/")
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
	updateUser := make(map[string]interface{})
	updateUser["_id"] = id
	if fUser.Name != nil {
		updateUser["name"] = *fUser.Name
	}
	if fUser.Age != nil {
		updateUser["age"] = *fUser.Age
	}
	if fUser.Email != nil {
		updateUser["email"] = *fUser.Email
	}
	fmt.Fprint(w, updateUser)
	//find and update user with id
	user, err := user.Update(h.Collection, updateUser)
	if err != nil {
		fmt.Fprint(w, errors.New("Unable to Update Item"))
		return
	}
	fmt.Fprint(w, user)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	//get id from url
	id := strings.TrimPrefix(r.URL.Path, "/api/user/")
	//find and delete user with id
	err := user.Delete(h.Collection, id)
	if err != nil {
		fmt.Fprint(w, errors.New("Unable to Delete Item"))
		return
	}
	fmt.Fprint(w, "successfully deleted")
}
