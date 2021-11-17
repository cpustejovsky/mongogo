package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

func ServerError(log *log.Logger, w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Error(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

type FormUser struct {
	Name   string
	Email  string
	Active bool
	Age    int
	Id     string
}

func DecodeUserForm(r *http.Request) (FormUser, error) {
	decoder := json.NewDecoder(r.Body)

	var user FormUser
	err := decoder.Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
