package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jpeach/githubfs/pkg/githubfs"
	"github.com/jpeach/githubfs/pkg/sysexits"
	"github.com/motemen/go-loghttp"
	flag "github.com/spf13/pflag"
)

// Progname ...
const Progname = "githubfs"

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... URL\n", Progname)
		fmt.Fprintf(os.Stderr, "Mount a GitHub repository as a filesystem\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
}

var traceFlag = flag.BoolP("http-trace", "t", false, "Enable HTTP tracing")
var _ = flag.BoolP("help", "?", false, "Print help")

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(sysexits.EX_USAGE)
	}

	client := &http.Client{}

	if *traceFlag {
		client.Transport = &loghttp.Transport{}
	}

	_, err := githubfs.New(
		githubfs.RepositoryOption(flag.Arg(0)),
		githubfs.ClientOption(client),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", Progname, err)
		os.Exit(sysexits.EX_SOFTWARE)
	}
}
