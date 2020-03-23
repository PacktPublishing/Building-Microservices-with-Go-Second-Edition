package handlers

import (
	"net/http"

	"github.com/hashicorp/go-hclog"
)

type Middleware struct {
	maxSize int64
	logger  hclog.Logger
}

// NewMiddleware creates a new middleware handlers
func NewMiddleware(maxContentLength int64, logger hclog.Logger) *Middleware {
	return &Middleware{maxSize: maxContentLength, logger: logger}
}

// CheckContentLengthMiddleware ensures that the content length is not greater than
// our allowed max.
// This can not be 100% depended on as Content-Length might be reported incorrectly
// however it is a fast first pass check.
func (mw *Middleware) CheckContentLengthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// if the content length is greater than the max reject the request
		if r.ContentLength > mw.maxSize {
			http.Error(rw, "Unable to save file, content too large", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
