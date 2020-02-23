package handlers

import (
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

	f.log.Info("handle standar post", "id", id, "filename", fn)

	// create a MaxBytesReader to ensure we do not read more than our max
	// over max size
	r.Body = http.MaxBytesReader(rw, r.Body, f.maxSize)

	f.saveFile(id, fn, rw, r.Body)
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

func (f *Files) saveFile(id, name string, rw http.ResponseWriter, r io.ReadCloser) {
	defer r.Close()
	r = http.MaxBytesReader(rw, r, f.maxSize)

	fp := filepath.Join(string(id), name)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
