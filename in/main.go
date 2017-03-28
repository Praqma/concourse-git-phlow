package main

import (
	"fmt"
	"os"
	"github.com/groenborg/pip/models"
	"encoding/json"
	"github.com/groenborg/pip/githandler"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

	destination := os.Args[1]

	err := os.MkdirAll(destination, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mkdir all fails:", err.Error())
		os.Exit(1)
	}

	var request models.InRequest

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "OS in parsing errored")
		os.Exit(1)
	}

	getRepo(request.Source.URL, destination)

	err = os.Chdir(destination)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not change dir:", err.Error())
		os.Exit(1)
	}

	out, err := githandler.Branch()
	if err != nil {
		fmt.Fprintln(os.Stderr, "branch failed:", err.Error())
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, out)

	_, err = githandler.PhlowReadyBranch()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Getting ready branch fail:", err.Error())
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "locating sha branch: %s \n", request.Version.Sha)

	fmt.Fprintf(os.Stderr, "Merging sha: %s with master\n", request.Version.Sha)
	err = githandler.Merge(strings.TrimSpace(request.Version.Sha))
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	GetMetadata(request.Version.Sha)
}

func GetMetadata(sha string) {
	ref, _ := githandler.CommitSha()
	author, _ := githandler.Author()
	date, _ := githandler.AuthorDate()

	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version: models.Version{Sha: sha},
		MetaData: models.Metadata{
			{"commit", ref},
			{"author", author},
			{"authordate", date},
		},
	})
}

func getRepo(url, path string) {
	fmt.Fprintf(os.Stderr, "Cloning into desitnation: %s from:  %s\n", path, url)
	_, err := githandler.Clone(url, path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "get repo failed:", err.Error())
		os.Exit(1)
	}
}
