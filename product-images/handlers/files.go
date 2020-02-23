package handlers

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-images/files"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

// Files is a handler for reading and writing files
type Files struct {
	log     hclog.Logger
	store   files.Storage
	maxSize int64
}

// NewFiles creates a new File handler
func NewFiles(s files.Storage, maxSize int64, l hclog.Logger) *Files {
	return &Files{store: s, maxSize: maxSize, log: l}
}

// ServeHTTP implements the http.Handler interface
func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// if the content length is greater than the max reject the request
	if r.ContentLength > f.maxSize {
		f.log.Error("Content length too large", "length", r.ContentLength)
		http.Error(rw, "Unable to save file, content too large", http.StatusBadRequest)
		return
	}

	// is this a multi-part form post
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		f.handleMultipart(rw, r)
		return
	}

	f.saveFile(rw, r)
}

// saveFile saves the contents of the request to a file from the request
func (f *Files) saveFile(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("handle standar post", "id", id, "filename", fn)

	// create a MaxBytesReader to ensure we do not read more than our max
	// over max size
	r.Body = http.MaxBytesReader(rw, r.Body, f.maxSize)

	fp := filepath.Join(id, fn)
	err := f.store.Save(fp, r.Body)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}

// handleMultipart handles multipart files
func (f *Files) handleMultipart(rw http.ResponseWriter, r *http.Request) {
	f.log.Info("Handle multipart POST")

	// Can be used to parse the multipart data but
	// will read the whole file
	// r.ParseForm()

	// not more efficient that ParseForm since the file has to be saved
	// to a temp file until the ID is known
	/*
		mr, err := r.MultipartReader()
		if err != nil {
			f.log.Error("Unable to read data", "error", err)
			http.Error(rw, "Unable to read multipart data", http.StatusBadRequest)
			return
		}

		for {
			np, err := mr.NextPart()
			if err == io.EOF {
				// done
				return
			}

			if err != nil {
				f.log.Error("Unable to read data", "error", err)
				http.Error(rw, "Unable to read multipart data", http.StatusBadRequest)
				return
			}

			switch np.FormName() {
			case "id":
				return
			case "file":
				// read the file

				return
			}

		}
	*/
}
