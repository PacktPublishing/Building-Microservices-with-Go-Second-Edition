package handlers

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GZipResponseMiddleware detects if the client can handle
// zipped content and if so returns the response in GZipped format
func GZipResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("ok done")

		// if client cant handle gzip send plain
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			//f.log.Debug("Unable to handle gzipped", "file", fp)
			next.ServeHTTP(rw, r)
			return
		}

		wr := NewWrappedResponseWriter(rw)

		// write the file
		next.ServeHTTP(wr, r)
	})
}

type WrappedResponseWriter struct {
	rw http.ResponseWriter
	gw io.Writer
}

func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {
	// client can handle gziped content send gzipped to speed up transfer
	// set the content encoding header for gzip
	rw.Header().Add("Content-Encoding", "gzip")

	// wrap the default writer in a gzip writer
	gw := gzip.NewWriter(rw)

	return &WrappedResponseWriter{rw, gw}
}

func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.rw.Header()
}

func (wr *WrappedResponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}

func (wr *WrappedResponseWriter) WriteHeader(statusCode int) {
	wr.rw.WriteHeader(statusCode)
}
