package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

var EmptyBodyError = errors.New("JSON body is empty")

func ServerError(log *log.Logger, w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Error(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//TODO: replace with generics
func DecodeForm(r *http.Request, form interface{}) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if len(body) <= 0 {
		return nil, EmptyBodyError
	}
	fmt.Println(body)
	err = json.Unmarshal(body, &form)
	if err != nil {
		return nil, err
	}
	return form, nil
}
