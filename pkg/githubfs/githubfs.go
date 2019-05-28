package githubfs

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-github/v25/github"
	"github.com/jpeach/githubfs/pkg/fs"
	"golang.org/x/exp/errors/fmt"
)

type repository struct {
	Client *github.Client
	HTTP   *http.Client
	URL    *url.URL

	Owner string
	Repo  string
	Ref   string
}

// Option is a githubfs repository client creation option.
type Option func(*repository) error

// ClientOption ...
func ClientOption(c *http.Client) Option {
	return func(r *repository) error {
		r.HTTP = c
		return nil
	}
}

// RepositoryOption ...
func RepositoryOption(repoURL string) Option {
	return func(r *repository) error {
		u, err := url.Parse(repoURL)
		if err != nil {
			return fmt.Errorf(
				"invalid repository URL '%s': %w", repoURL, err)
		}

		r.URL = u
		return nil
	}
}

// New ...
func New(options ...Option) (fs.Repository, error) {
	r := &repository{
		Ref: "master",
	}

	for _, o := range options {
		if err := o(r); err != nil {
			return nil, err
		}
	}

	// TODO(jpeach): Make sure that this actually works for GHE.
	r.Client, _ = github.NewEnterpriseClient(
		fmt.Sprintf("https://api.%s/", r.URL.Hostname()),
		fmt.Sprintf("https://api.%s/", r.URL.Hostname()),
		r.HTTP)

	r.Owner = strings.TrimLeft(path.Dir(r.URL.Path), "/")
	r.Repo = strings.TrimSuffix(strings.TrimLeft(path.Base(r.URL.Path), "/"), ".git")

	return r, nil
}

// ReadDir reads a directory from the GitHub repository. Note that
// this uses the GitHub GetContents API, which is documented to only
// return up to 1000 entries.
//
// See https://developer.github.com/v3/repos/contents/#get-contents.
func (r *repository) ReadDir(dirname string) ([]os.FileInfo, error) {
	file, dirent, resp, err := r.Client.Repositories.GetContents(
		context.TODO(),
		r.Owner,
		r.Repo,
		path.Clean(dirname),
		&github.RepositoryContentGetOptions{
			Ref: r.Ref,
		})
	if err != nil {
		return nil, wrapResponseError(resp, err)
	}

	// The path is a file, not a directory.
	if file != nil {
		return nil, fmt.Errorf("failed to read %s: %w",
			dirname, syscall.ENOTDIR)
	}

	// It's not a path, not an error, and not a directory. Who knows?
	if dirent == nil {
		return nil, fmt.Errorf("failed to read %s: %w",
			dirname, syscall.EIO)
	}

	info := make([]os.FileInfo, 0, len(dirent))
	for _, entry := range dirent {
		info = append(info, FileInfoFromContent(entry))
	}

	return info, nil
}

func (r *repository) ReadFile(filename string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *repository) Lookup(pathname string) (os.FileInfo, error) {
	file, dirent, resp, err := r.Client.Repositories.GetContents(
		context.TODO(),
		r.Owner,
		r.Repo,
		path.Clean(pathname),
		&github.RepositoryContentGetOptions{
			Ref: r.Ref,
		})
	if err != nil {
		return nil, wrapResponseError(resp, err)
	}

	if file != nil {
		return FileInfoFromContent(file), nil
	}

	if dirent != nil {
		return &FileInfo{
			name:  path.Base(pathname),
			mtime: time.Now(),
			ftype: FileTypeDirectory,
		}, nil
	}

	return nil, fmt.Errorf("failed to lookup '%s': %w",
		pathname, syscall.EIO)
}
