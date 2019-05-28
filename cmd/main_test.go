package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jpeach/githubfs/pkg/githubfs"
	"github.com/motemen/go-loghttp"
)

func TestGitHubLookup(t *testing.T) {
	url := "https://github.com/jpeach/vimrc.git"
	repo, err := githubfs.New(
		githubfs.RepositoryOption(url),
		githubfs.ClientOption(&http.Client{
			Transport: &loghttp.Transport{},
		}),
	)

	if err != nil {
		t.Fatalf("failed to init '%s': %s", url, err)
	}

	info, err := repo.Lookup("/doc/ack.txt")
	if err != nil {
		t.Fatalf("lookup failed: %s", err)
	}

	fmt.Printf("%10d %s\n", info.Size(), info.Name())
}

func TestGitHubReadDir(t *testing.T) {
	url := "https://github.com/jpeach/vimrc.git"
	repo, err := githubfs.New(
		githubfs.RepositoryOption(url),
		githubfs.ClientOption(&http.Client{
			Transport: &loghttp.Transport{},
		}),
	)

	if err != nil {
		t.Fatalf("failed to init '%s': %s", url, err)
	}

	dirent, err := repo.ReadDir("/")
	if err != nil {
		t.Fatalf("failed to readdir '%s': %s", "/", err)
	}

	for _, d := range dirent {
		fmt.Printf("%10d %s\n", d.Size(), d.Name())
	}

	dirent, err = repo.ReadDir("/doc")
	if err != nil {
		t.Fatalf("failed to readdir '%s': %s", "/", err)
	}

	for _, d := range dirent {
		fmt.Printf("%10d %s\n", d.Size(), d.Name())
	}

}
