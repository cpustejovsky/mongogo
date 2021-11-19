package helpers_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cpustejovsky/mongogo/handlers"
	"github.com/cpustejovsky/mongogo/helpers"
)

var stubUserHandlers = handlers.Handler{}

func TestDecodeUserForm(t *testing.T) {
	fullJSONBody := []byte(`{
		"name": "Charles",
		"email": "charles.pustejovsky@gmail.com",
		"age": 28,
		"active": true	}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(fullJSONBody))
	req.Header.Set("Content-Type", "application/json")
	_, err := helpers.DecodeUserForm(req)
	if err != nil {
		t.Errorf("got an error:\n%v", err)
	}
	t.Run("Returns emptyBodyError for empty body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(``)))
		req.Header.Set("Content-Type", "application/json")
		_, err := helpers.DecodeUserForm(req)
		if err != helpers.EmptyBodyError {
			t.Errorf(`got "%v", wanted "%v"`, err, helpers.EmptyBodyError)
		}
	})
}

func TestMissingPropertyError(t *testing.T) {
	errStr := fmt.Sprintf(helpers.MissingPropertyErrorTemplateString, "foo, bar, baz")
	want := errors.New(errStr)
	missingProps := []string{"foo", "bar", "baz"}
	got := helpers.MissingPropertyError(missingProps)
	if got.Error() != want.Error() {
		t.Errorf("\nwant:\n%v\ngot:\n%v", want, got)
	}
}
