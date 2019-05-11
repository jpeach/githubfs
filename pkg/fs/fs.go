package fs

// Repository ...
type Repository interface {
	Readdir(path string) error
	Stat(path string) error
}
