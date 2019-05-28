package githubfs

import (
	"syscall"

	"github.com/google/go-github/v25/github"
	"golang.org/x/exp/errors/fmt"
)

// wrapResponseError converts a GitHub error response into a standard error.
func wrapResponseError(response *github.Response, err error) error {
	errno := syscall.EIO

	switch response.StatusCode {
	case 403:
		errno = syscall.EPERM
	case 404:
		errno = syscall.ENOENT
	}

	return fmt.Errorf("%s: %w", err, errno)
}
