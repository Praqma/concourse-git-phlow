package main

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/praqma/concourse-git-phlow/githandler"
	"github.com/praqma/concourse-git-phlow/models"
	"github.com/praqma/concourse-git-phlow/repo"
	"github.com/praqma/git-phlow/phlow"
	"strings"
)

//Strategy ...
type Strategy interface {
	Checkout(string) error
	MergeFF(string) error
	RebaseOnto(string) error
}

//GitStrategy ...
type GitStrategy struct {
}

func (g *GitStrategy) Checkout(br string) error {
	return githandler.CheckOut(br)
}

func (g *GitStrategy) MergeFF(br string) error {
	return githandler.MergeFFO(br)
}

func (g *GitStrategy) RebaseOnto(br string) error {
	return githandler.RebaseOnto(br)
}

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

	//Check if sha from in is the same as master sha
	//git branch master --contains <commit>
	cco, err := githandler.ContainsCommit(request.Source.Master, request.Version.Sha)
	fmt.Fprintln(os.Stderr, "cco: "+cco)
	if strings.TrimSpace(cco) != "" {
		//Exists on master
		fmt.Fprintln(os.Stderr, "Found sha on master, no need for integration")
		repo.WriteRDYBranch("")
		SendMetadata(request.Version.Sha)
		os.Exit(0)
	}

	//list all branches
	out := githandler.BranchList()
	fmt.Fprintln(os.Stderr, out)

	//retrieve the next ready branch from origin with prefix
	rbn := phlow.UpNext("origin", request.Source.PrefixReady)
	fmt.Fprintln(os.Stderr, "Ready branch found: "+rbn)
	if rbn == "" {
		fmt.Fprintf(os.Stderr, "no branches with: %s available for integration with: %s \n", request.Source.PrefixReady, request.Source.Master)
		fmt.Fprintln(os.Stderr, "Exiting build")
		repo.WriteRDYBranch("") //write an empty name

		SendMetadata(request.Version.Sha)
		os.Exit(0)
	}

	//Checkout ready branch to get a local copy
	if err = githandler.CheckOut(rbn); err != nil {
		fmt.Fprintln(os.Stderr, "checkout failed: ", err.Error())
		os.Exit(1)
	}

	//Checkout master
	if err = githandler.CheckOut(request.Source.Master); err != nil {
		fmt.Fprintln(os.Stderr, "could not checkout main branch:", err.Error())
		os.Exit(1)
	}

	//Names the branch to the name plus wip prefix
	wipBranchName := request.Source.PrefixWip + strings.TrimPrefix(rbn, request.Source.PrefixReady)
	u := repo.FormatURL(request.Source.URL, request.Source.Username, request.Source.Password)
	fmt.Fprintln(os.Stderr, "wip branch: "+wipBranchName)
	fmt.Fprintln(os.Stderr, "ready branch: "+rbn)
	repo.RenameRemoteBranch(u, wipBranchName, rbn)

	strategy := GitStrategy{}
	err = ApplyAndRunStrategy(request.Source.Master, rbn, &strategy)
	if err != nil {
		err := githandler.CheckOut(wipBranchName)
		fmt.Fprintln(os.Stderr, err)

		repo.RenameRemoteBranch(u, "failed/"+rbn, wipBranchName)
		fmt.Fprintln(os.Stderr, "Merge failed, Aborting integration")
		os.Exit(1)
	}

	repo.WriteRDYBranch(wipBranchName)
	fmt.Fprintf(os.Stderr, "saving wip branch name: %s to dest: %s for use in out", wipBranchName, ".git/git-phlow-ready-branch")
	SendMetadata(request.Version.Sha)
}

//ApplyStrategy
// 0 - ff-only merge
// 1 - try rebase
func ApplyAndRunStrategy(master string, ready string, s Strategy) (err error) {

	var rb = func() error {
		//checkout ready before rebase
		if err := s.Checkout(ready); err != nil {
			fmt.Fprintf(os.Stderr, "could not checkout %s \n", ready)
			return err
		}

		if err := s.RebaseOnto(master); err != nil {
			fmt.Fprintln(os.Stderr, "not able to rebase")
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		return nil
	}

	var ff = func() error {
		//checkout master before fast-forward merge
		if err := s.Checkout(master); err != nil {
			fmt.Fprintf(os.Stderr, "Could not checkout %s \n", master)
			return err
		}
		if err := s.MergeFF(ready); err != nil {
			fmt.Fprintln(os.Stderr, "not able to fast forward")
			return err
		}
		return nil
	}

	if err = ff(); err == nil {
		fmt.Fprintln(os.Stderr, "fast-forward success")
		return nil
	} else {
		if err = rb(); err != nil {
			fmt.Fprintln(os.Stderr, "rebase fail")
			return err
		} else {
			if err = ff(); err != nil {
				fmt.Fprintln(os.Stderr, "fast-forward fail")
				return err
			}
			fmt.Fprintln(os.Stderr, "fast-forward success after rebases")
			return nil
		}
	}
}

//GetMetadata ...
//sends the metadata to output
func SendMetadata(sha string) {
	ref, _ := githandler.CommitSha()
	author, _ := githandler.Author()
	date, _ := githandler.AuthorDate()

	str, err := json.Marshal(models.InResponse{
		Version: models.Version{Sha: sha},
		MetaData: models.Metadata{
			{"commit", ref},
			{"author", author},
			{"authordate", date},
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Fprint(os.Stdout, string(str))
	//json.NewEncoder(os.Stdout).Encode()
}
