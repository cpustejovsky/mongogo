package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/cpustejovsky/mongogo/helpers"
)

type Middleware struct {
	Logger *zap.SugaredLogger
}

func (m *Middleware) SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		// Any code here will execute on the way down the chain.
		next.ServeHTTP(w, r)
		// Any code here will execute on the way back up the chain.
	})
}

//TODO: add type safety around requestId key and value
func (m *Middleware) SetRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "requestId", id.String())))
	})
}

func (m *Middleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value("requestId")
		m.Logger.Infow("Request", "Remote Address", r.RemoteAddr, "Proto", r.Proto, "Method", r.Method, "URI", r.URL.RequestURI(), "ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "ID", id)))
	})
}

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				id := r.Context().Value("requestId")
				m.Logger.Errorw("Error", "URI", r.URL.RequestURI(), "ID", id)
				w.Header().Set("Connection", "close")
				helpers.ServerError(m.Logger, w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
