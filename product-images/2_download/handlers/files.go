package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-images/1_restful/files"
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

func (f *Files) saveFile(id, name string, rw http.ResponseWriter, r io.ReadCloser) error {
	// MaxBytesReader ensures that a file does not exceed the given size
	r = http.MaxBytesReader(rw, r, f.maxSize)
	defer r.Close()

	// create the file path and save the file
	fp := filepath.Join(string(id), name)
	err := f.store.Save(fp, r)
	err = nil
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}

	return err
}
