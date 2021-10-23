package helpers

import (
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
