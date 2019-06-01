package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jpeach/githubfs/pkg/fs"
	"github.com/jpeach/githubfs/pkg/githubfs"
	"github.com/motemen/go-loghttp"
)

const url = "https://github.com/jpeach/vimrc.git"

func newRepo(url string) (fs.Repository, error) {
	return githubfs.New(
		githubfs.RepositoryOption(url),
		githubfs.ClientOption(&http.Client{
			Transport: &loghttp.Transport{},
		}),
	)
}

func TestGitHubLookup(t *testing.T) {
	repo, err := newRepo(url)
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
	repo, err := newRepo(url)
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

func TestGitHubReadTree(t *testing.T) {
	repo, err := newRepo(url)
	if err != nil {
		t.Fatalf("failed to init '%s': %s", url, err)
	}

	g := repo.(*githubfs.Repository)
	tree, err := g.ReadTree()
	if err != nil {
		t.Fatalf("failed to read tree: %s", err)
	}

	fmt.Printf("SHA: %s\n", tree.GetSHA())
	fmt.Printf("Truncated: %t\n", tree.GetTruncated())

	for _, e := range tree.Entries {
		fmt.Printf("\n%s\n", e)
	}
}
