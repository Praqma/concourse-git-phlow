package main

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/praqma/concourse-git-phlow/githandler"
	"github.com/praqma/concourse-git-phlow/models"
	"github.com/praqma/concourse-git-phlow/repo"
	"github.com/praqma/git-phlow/phlow"
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

	githandler.Status()

	rbn := phlow.UpNext("origin", request.Source.PrefixReady)
	if rbn == "" {
		fmt.Fprintln(os.Stderr, "No ready branch to integrate with master.. Exiting build")
		repo.WriteRDYBranch("") //write an empty name
		SendMetadata(request.Version.Sha)
		os.Exit(0)
	}

	err = githandler.CheckOut(rbn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "checkout failed: ", err.Error())
		os.Exit(1)
	}

	err = githandler.CheckOut(request.Source.Master)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not checkout main branch:", err.Error())
		os.Exit(1)
	}

	//Names the branch to the name plus wip prefix
	wipBranchName := request.Source.PrefixWip + rbn
	u := repo.FormatURL(request.Source.URL, request.Source.Username, request.Source.Password)
	fmt.Fprintln(os.Stderr, "wip branch: "+wipBranchName)
	fmt.Fprintln(os.Stderr, "ready branch: "+rbn)

	repo.RenameRemoteBranch(u, wipBranchName, rbn)

	fmt.Fprintf(os.Stderr, "Merging sha: %s with master\n", request.Version.Sha)
	err = githandler.Merge(request.Version.Sha)
	if err != nil {
		githandler.CheckOut(wipBranchName) // ARHGGSJAHGDHJSAD
		repo.RenameRemoteBranch(u, "failed/"+rbn, wipBranchName)
		fmt.Fprintln(os.Stderr, "Merge failed, Aborting integration")
		os.Exit(1)
	}

	repo.WriteRDYBranch(wipBranchName)
	SendMetadata(request.Version.Sha)
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
