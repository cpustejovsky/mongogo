package helpers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cpustejovsky/mongogo/handlers"
	"github.com/cpustejovsky/mongogo/helpers"
	"github.com/cpustejovsky/mongogo/internal/models"
)

var stubUserHandlers = handlers.Handler{}
var user models.User

func TestDecodeForm(t *testing.T) {
	fullJSONBody := []byte(`{
		"name": "Charles",
		"email": "charles.pustejovsky@gmail.com",
		"age": 28,
		"active": true	}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(fullJSONBody))
	req.Header.Set("Content-Type", "application/json")
	_, err := helpers.DecodeForm(req, user)
	if err != nil {
		t.Errorf("got an error:\n%v", err)
	}
	t.Run("Returns emptyBodyError for empty body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(``)))
		req.Header.Set("Content-Type", "application/json")
		_, err := helpers.DecodeForm(req, user)
		if err != helpers.EmptyBodyError {
			t.Errorf(`got "%v", wanted "%v"`, err, helpers.EmptyBodyError)
		}
	})
}
