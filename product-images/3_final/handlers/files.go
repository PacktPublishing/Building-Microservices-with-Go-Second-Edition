package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

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

// SaveFileREST saves the contents of the request to a file from the request
func (f *Files) SaveFileREST(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle RESTful POST", "id", id, "filename", fn)

	// create a MaxBytesReader to ensure we do not read more than our max
	// over max size
	r.Body = http.MaxBytesReader(rw, r.Body, f.maxSize)

	err := f.saveFile(id, fn, rw, r.Body)
	if err != nil {
		return
	}

	// write the json response
	resp := map[string]string{}
	resp["status"] = "200"
	resp["message"] = "ok"

	json.NewEncoder(rw).Encode(resp)
}

// SaveMultipart handles multipart files
func (f *Files) SaveMultipart(rw http.ResponseWriter, r *http.Request) {
	f.log.Info("Handle multipart POST")
	err := r.ParseMultipartForm(f.maxSize)
	if err != nil {
		f.log.Error("Unable to parse form", "error", err)
		http.Error(rw, "Unable to read multipart data", http.StatusBadRequest)
		return
	}

	// check we have an id and a file
	id, _ := strconv.Atoi(r.FormValue("id"))
	ff, fh, err := r.FormFile("file")
	if id < 1 || err != nil {
		f.log.Error("Unable to validate form")
		http.Error(rw, "Please provide a valid id and file to upload", http.StatusPreconditionFailed)
		return
	}

	// save the file
	f.saveFile(r.FormValue("id"), fh.Filename, rw, ff)
}

func (f *Files) saveFile(id, name string, rw http.ResponseWriter, r io.ReadCloser) error {
	defer r.Close()
	r = http.MaxBytesReader(rw, r, f.maxSize)

	fp := filepath.Join(string(id), name)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}

	return err
}
