package files

import "io"

import "os"

// Storage defines the behavior for file operations
// Implementations may be of the time local disk, or cloud storage, etc
type Storage interface {
	Save(path string, file io.Reader) error
	Get(path string) (*os.File, error)
}
