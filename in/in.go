package main

import (
	"fmt"
	"os"
	"github.com/groenborg/pip/models"
	"encoding/json"
	"github.com/groenborg/pip/githandler"
	"github.com/groenborg/pip/repo"
	"io/ioutil"
)

func main() {

	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

	destination := os.Args[1]

	err := os.MkdirAll(destination, 0755)
	repo.Check(err, "cannot make dir")

	var request models.InRequest

	err = json.NewDecoder(os.Stdin).Decode(&request)
	repo.Check(err, "OS in parsing errored")

	repo.CloneRepoSource(request.Source.URL, destination, request.Source.Username, request.Source.Password)

	err = os.Chdir(destination)
	repo.Check(err, "could not change dir")

	rbn, err := githandler.PhlowReadyBranch()
	repo.Check(err, "an error")


	err = ioutil.WriteFile(".git/git-phlow-ready-branch", []byte(rbn), 0655)
	repo.Check(err, "could not write to file")

	if rbn == "" {
		fmt.Fprintln(os.Stderr, "No ready branch to integrate with master.. Exiting build")
		GetMetadata(request.Version.Sha)
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "locating sha branch: %s \n", request.Version.Sha)
	fmt.Fprintf(os.Stderr, "Merging sha: %s with master\n", request.Version.Sha)

	err = githandler.Merge(request.Version.Sha)
	repo.Check(err, "could not merge")

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
