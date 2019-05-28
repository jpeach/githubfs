package fs

import "os"

// Repository ...
type Repository interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	Lookup(path string) (os.FileInfo, error)
}
