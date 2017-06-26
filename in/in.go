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
	"github.com/praqma/concourse-git-phlow/concourse"
	"log"
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

//CheckOut ..
func (g *GitStrategy) Checkout(br string) error {
	return githandler.CheckOut(br)
}

//MergeFFO ...
func (g *GitStrategy) MergeFF(br string) error {
	return githandler.MergeFFO(br)
}

//RebaseOnto ...
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
	if err != nil {
		log.Panicln(err)
	}

	var request models.InRequest

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Fprintln(os.Stderr, "Resource Version " + repo.Version)


	repo.CloneRepoSource(request.Source.URL, destination, request.Source.Username, request.Source.Password)

	err = os.Chdir(destination)
	if err != nil {
		log.Panicln(err)
	}

	RunPhlow(&request)
}

//RunPhlow ...
//Runs the phlow workflow
func RunPhlow(request *models.InRequest) {

	//Verify if the commit from IN is already on master
	cco, err := githandler.ContainsCommit(request.Source.MainBranch, request.Version.Sha)
	if err != nil {
		log.Panicln(err)
	}

	if strings.TrimSpace(cco) != "" {
		fmt.Fprintln(os.Stderr, "Found sha on master, no need for integration")
		fmt.Fprintln(os.Stderr, cco)
		repo.WriteRDYBranch("")
		concourse.SendMetadata(request.Version.Sha)
		os.Exit(0)
	}

	gitBranches, err := githandler.BranchList()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Fprintln(os.Stderr, gitBranches)

	//retrieve the next ready branch from origin with prefix
	rbn := phlow.UpNext("origin", request.Source.PrefixReady)
	if rbn == "" {
		fmt.Fprintf(os.Stderr, "no branches with: %s available for integration with: %s .. Exiting\n", request.Source.PrefixReady, request.Source.MainBranch)
		repo.WriteRDYBranch("") //write an empty name
		concourse.SendMetadata(request.Version.Sha)
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, "Target branch: %s", rbn)

	//Checkout ready branch to get a local copy
	if err = githandler.CheckOut(rbn); err != nil {
		log.Panicln(err)
	}

	//Checkout master
	if err = githandler.CheckOut(request.Source.MainBranch); err != nil {
		log.Panicln(err)
	}

	//Names the branch to the name plus wip prefix
	wbn := request.Source.PrefixWip + strings.TrimPrefix(rbn, request.Source.PrefixReady)
	url := repo.FormatURL(request.Source.URL, request.Source.Username, request.Source.Password)
	repo.RenameRemoteBranch(url, wbn, rbn)

	strategy := GitStrategy{}
	err = ApplyAndRunStrategy(request.Source.MainBranch, rbn, &strategy)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Merge failed, Aborting integration")
		os.Exit(1)
	}

	repo.WriteRDYBranch(wbn)
	concourse.SendMetadata(request.Version.Sha)
}

//ApplyAndRunStrategy ...
//Run phlow pretested-integration strategy
func ApplyAndRunStrategy(master string, ready string, s Strategy) (err error) {

	var rb = func() error {
		//checkout ready before rebase
		if err := s.Checkout(ready); err != nil {
			fmt.Fprintf(os.Stderr, "could not checkout %s \n", ready)
			return err
		}
		//rebase
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
			fmt.Fprintln(os.Stderr, " not able to fast forward")
			return err
		}
		return nil
	}

	if err = ff(); err == nil {
		//First try ff-merge
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
