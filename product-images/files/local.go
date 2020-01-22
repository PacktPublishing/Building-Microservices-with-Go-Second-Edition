package files

import (
	"io"
	"os"
	"path/filepath"
)

// Local is an implementation of the Storage interface which works with the
// local disk on the current machine
type Local struct {
	basePath string
}

// NewLocal creates a new Local filesytem with the given base path
func NewLocal(basePath string) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Local{p}, nil
}

// Save the contents of the Writer to the given path
func (l *Local) Save(path string, contents io.Reader) error {
	// append the given path to the base path
	fp := filepath.Join(l.basePath, path)

	// create a new file at the path
	f, err := os.Create(fp)
	if err != nil {
		return err
	}

	// write the contents to the new file
	_, err = io.Copy(f, contents)

	return err
}

// Get the file at the given path and return a Reader
func (l *Local) Get(path string) (io.Reader, error) {
	return nil, nil
}
