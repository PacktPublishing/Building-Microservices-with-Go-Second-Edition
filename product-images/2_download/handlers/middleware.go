package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

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

// GZipResponseMiddleware detects if the client can handle
// zipped content and if so returns the response in GZipped format
func (mw *Middleware) GZipResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// if client cant handle gzip send plain
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			//f.log.Debug("Unable to handle gzipped", "file", fp)
			mw.logger.Debug("Handle Gziped file")

			// client can handle gziped content send gzipped to speed up transfer
			// set the content encoding header for gzip
			rw.Header().Add("Content-Encoding", "gzip")

			// file server sets the content stream
			// nice
			//rw.Header().Add("Content-Type", "application/octet-stream")

			rw = NewWrappedResponseWriter(rw)
			defer rw.(http.Flusher).Flush()
		}

		// write the file
		next.ServeHTTP(rw, r)
	})
}

// WrappedResponseWriter wrapps the default http.ResponseWriter in a GZip stream
type WrappedResponseWriter struct {
	http.ResponseWriter
	gw io.Writer
}

// NewWrappedResponseWriter returns a new wrapped response writer
func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {
	// wrap the default writer in a gzip writer
	gw := gzip.NewWriter(rw)

	return &WrappedResponseWriter{rw, gw}
}

// Write overrides the http.ResponseWriter Write method and ensures
// that data is written throught he gzip writer not direct to the output stream
func (wr *WrappedResponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}

// Flush implements the Flusher interface which allows ResponseWriters
func (wr *WrappedResponseWriter) Flush() {
	wr.gw.(http.Flusher).Flush()
}
