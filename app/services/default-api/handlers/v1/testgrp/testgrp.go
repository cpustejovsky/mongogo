package testgrp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cpustejovsky/mongogo/foundation/web"
	"github.com/cpustejovsky/mongogo/helpers"
	"github.com/cpustejovsky/mongogo/internal/models/mongodb/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Handlers manages the set of check enpoints.
type Handlers struct {
	Logger     *zap.SugaredLogger
	Collection *mongo.Collection
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "Howdy!",
	}
	statusCode := http.StatusOK
	return web.Respond(ctx, w, status, statusCode)
}

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *Handlers) PingPanic(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("requestId")
	idstr := fmt.Sprintf("Request ID: %v\n", id)
	w.Write([]byte(idstr))
	panic("foo")
}

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprint(w, errors.New("Unable to Create Item"))
		return
	}
	ju, err := json.Marshal(user)
	if err != nil {
		h.Logger.Error(err)
		fmt.Fprint(w, errors.New("Could not marshall object"))
	}
	fmt.Fprint(w, string(ju))
}

func (h *Handlers) Fetch(w http.ResponseWriter, r *http.Request) {
	//get id from url
	id := strings.TrimPrefix(r.URL.Path, "/api/user/")
	//find user by id and return
	user, err := user.Fetch(h.Collection, id)
	if err != nil {
		h.Logger.Error(err)
		fmt.Fprint(w, errors.New("Unable to Fetch Item"))
		return
	}
	ju, err := json.Marshal(user)
	if err != nil {
		h.Logger.Error(err)
		fmt.Fprint(w, errors.New("Could not marshall object"))
	}
	fmt.Fprint(w, string(ju))
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//get id from url
	id := strings.TrimPrefix(r.URL.Path, "/api/user/")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helpers.ServerError(h.Logger, w, err)
		return
	}
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
	updateUser["_id"] = oid
	if fUser.Name != nil {
		updateUser["name"] = *fUser.Name
	}
	if fUser.Age != nil {
		updateUser["age"] = *fUser.Age
	}
	if fUser.Email != nil {
		updateUser["email"] = *fUser.Email
	}
	//find and update user with id
	user, err := user.Update(h.Collection, updateUser)
	if err != nil {
		fmt.Fprint(w, errors.New("Unable to Update Item"))
		return
	}
	ju, err := json.Marshal(user)
	if err != nil {
		h.Logger.Error(err)
		fmt.Fprint(w, errors.New("Could not marshall object"))
	}
	fmt.Fprint(w, string(ju))
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
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
