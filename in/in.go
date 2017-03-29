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
	repo.Check(err, "error in locating readybranch")

	if rbn == "" {
		fmt.Fprintln(os.Stderr, "No ready branch to integrate with master.. Exiting build")
		WriteRDYBranch("") //write an empty name
		SendMetadata(request.Version.Sha)
		os.Exit(0)
	}

	//Names the branch to the name plus wip prefix
	wipBranchName := request.Source.PrefixWip + rbn
	u := repo.FormatURL(request.Source.URL, request.Source.URL, request.Source.Password)

	fmt.Fprintf(os.Stderr, "Merging sha: %s with master\n", request.Version.Sha)
	err = githandler.Merge(request.Version.Sha)
	if err != nil {
		repo.RenameRemoteBranch(u, "failed/"+rbn, rbn)
		os.Exit(1)
	}

	repo.RenameRemoteBranch(u, wipBranchName, rbn)
	WriteRDYBranch(wipBranchName)
	SendMetadata(request.Version.Sha)
}

//WriteRDYBranch ...
//writes the name of the branch to the file
func WriteRDYBranch(name string) {
	err := ioutil.WriteFile(".git/git-phlow-ready-branch", []byte(name), 0655)
	if err != nil {
		os.Exit(1)
	}
}

//GetMetadata ...
//sends the metadata to output
func SendMetadata(sha string) {
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
