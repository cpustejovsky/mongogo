package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"

	"github.com/cpustejovsky/mongogo/internal/models"
	"go.uber.org/zap"
)

var EmptyBodyError = errors.New("JSON body is empty")

func ServerError(log *zap.SugaredLogger, w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Error(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//TODO: replace with generics
func DecodeUserForm(r *http.Request) (models.FormUser, error) {
	var form models.FormUser
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return form, err
	}
	if len(body) <= 0 {
		return form, EmptyBodyError
	}
	err = json.Unmarshal(body, &form)
	if err != nil {
		return form, err
	}
	return form, nil
}

var MissingPropertyErrorTemplateString = "Missing the following property or properties: %v"

func MissingPropertyError(props []string) error {
	var propStr string
	for i, prop := range props {
		if i > 0 {
			propStr += ", "
		}
		propStr += prop
	}
	return errors.New(fmt.Sprintf(MissingPropertyErrorTemplateString, propStr))
}
