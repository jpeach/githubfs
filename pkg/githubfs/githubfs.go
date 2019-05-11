package githubfs

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strings"

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

func (r *repository) Readdir(path string) error {
	tree, _, err := r.Client.Git.GetTree(
		context.TODO(),
		r.Owner,
		r.Repo,
		r.Ref,
		false /* recursive */)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", tree)
	return nil
}

func (r *repository) Stat(path string) error {
	return fmt.Errorf("not implemented")
}
