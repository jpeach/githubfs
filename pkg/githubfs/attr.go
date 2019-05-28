package githubfs

import (
	"os"
	"time"

	"github.com/google/go-github/v25/github"
)

// FileType ...
type FileType int

const (
	// FileTypeFile ...
	FileTypeFile FileType = iota

	// FileTypeDirectory ...
	FileTypeDirectory

	// FileTypeSymlink ...
	FileTypeSymlink
)

// FileInfo implements os.FileInfo
type FileInfo struct {
	name  string
	size  int64
	ftype FileType
	mode  os.FileMode
	mtime time.Time
}

// FileInfoFromContent ...
func FileInfoFromContent(content *github.RepositoryContent) *FileInfo {
	f := &FileInfo{}

	f.name = content.GetName()
	f.size = int64(content.GetSize())
	f.mtime = time.Now()

	// XXX(jpeach); the Contents API doesn't give us
	// the file mode, but the Tree API does. However,
	// with the Trees API, there's no way to get a
	// partial tree, which means we have to take a bit
	// download up front.

	switch *content.Type {
	case "file":
		f.mode = 0644
		f.ftype = FileTypeFile
	case "dir":
		f.mode = 0755
		f.ftype = FileTypeDirectory
	case "symlink":
		f.mode = 0644
		f.ftype = FileTypeSymlink

	}

	return f
}

// Name ...
func (f *FileInfo) Name() string {
	return f.name
}

// Size ...
func (f *FileInfo) Size() int64 {
	return f.size
}

// Mode ...
func (f *FileInfo) Mode() os.FileMode {
	return f.mode
}

// ModTime ...
func (f *FileInfo) ModTime() time.Time {
	return f.mtime
}

// IsDir ...
func (f *FileInfo) IsDir() bool {
	return f.ftype == FileTypeDirectory
}

// Sys ...
func (f *FileInfo) Sys() interface{} {
	return nil
}
