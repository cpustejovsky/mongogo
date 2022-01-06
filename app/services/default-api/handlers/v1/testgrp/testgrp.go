package testgrp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
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

func (h *Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "Howdy!",
	}
	statusCode := http.StatusOK
	return web.Respond(ctx, w, status, statusCode)
}

func (h *Handlers) Ping(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return web.Respond(ctx, w, []byte("OK"), 200)
}

func (h *Handlers) TestPanic(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.Context().Value("requestId")
	idstr := fmt.Sprintf("Request ID: %v\n", id)
	if n := rand.Intn(100); n%2 == 0 {
		panic("testing panic!")
	}
	return web.Respond(ctx, w, idstr, 200)
}

func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//get JSON body and decode
	fUser, err := helpers.DecodeUserForm(r)
	if err != nil {
		if err == helpers.EmptyBodyError {
			fmt.Fprint(w, err)
			return web.Respond(ctx, w, err, 400)
		}
		helpers.ServerError(h.Logger, w, err)
		return web.Respond(ctx, w, err, 400)
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
		return web.Respond(ctx, w, helpers.MissingPropertyError(missingProperties), 400)
	}
	//create new document within mongodb table
	user, err := user.Create(h.Collection, fUser)
	if err != nil {
		return web.Respond(ctx, w, errors.New("unable to create item"), 400)
	}
	ju, err := json.Marshal(user)
	if err != nil {
		h.Logger.Error(err)
		fmt.Fprint(w, errors.New("could not marshall object"))
	}
	return web.Respond(ctx, w, string(ju), 200)
}

func (h *Handlers) Fetch(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//get id from url
	id := strings.TrimPrefix(r.URL.Path, "/api/user/")
	//find user by id and return
	user, err := user.Fetch(h.Collection, id)
	if err != nil {
		h.Logger.Error(err)
		return web.Respond(ctx, w, errors.New("unable to fetch item"), 400)
	}
	ju, err := json.Marshal(user)
	if err != nil {
		h.Logger.Error(err)
		return web.Respond(ctx, w, errors.New("could not marshall object"), 400)
	}
	return web.Respond(ctx, w, string(ju), 200)
}

func (h *Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	//get id from url
	id := strings.TrimPrefix(r.URL.Path, "/api/user/")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return web.Respond(ctx, w, err, 400)
	}
	//get JSON body and decode
	fUser, err := helpers.DecodeUserForm(r)
	if err != nil {
		if err == helpers.EmptyBodyError {
			return web.Respond(ctx, w, err, 400)
		}
		return web.Respond(ctx, w, err, 400)
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
		return web.Respond(ctx, w, errors.New("unable to update item"), 400)
	}
	ju, err := json.Marshal(user)
	if err != nil {
		h.Logger.Error(err)
		return web.Respond(ctx, w, errors.New("could not marshall object"), 400)
	}
	return web.Respond(ctx, w, string(ju), 200)
}

func (h *Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//get id from url
	id := strings.TrimPrefix(r.URL.Path, "/api/user/")
	//find and delete user with id
	err := user.Delete(h.Collection, id)
	if err != nil {
		return web.Respond(ctx, w, errors.New("unable to celete item"), 400)
	}
	return web.Respond(ctx, w, "successfully deleted", 200)
}
