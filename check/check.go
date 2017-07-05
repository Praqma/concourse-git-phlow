package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/praqma/concourse-git-phlow/githandler"
	"github.com/praqma/concourse-git-phlow/models"
	"github.com/praqma/concourse-git-phlow/repo"
	"github.com/praqma/git-phlow/phlow"
	"log"
)

func main() {
	var request models.CheckRequest
	var ref string
	cacheDir := "/cache"

	destination := os.Getenv("TMPDIR")
	if destination == "" {
		log.Panicln("TMPDIR Missing: ", destination)
	}

	destination = destination + cacheDir

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		log.Panicln("Unable to parse json input", err)
		os.Exit(1)
	}

	//validating if the repository is cached
	if doesExist(destination) {
		ref = getRef(destination, request)
	} else {
		repo.CloneRepoSource(request.Source.URL, destination, request.Source.Username, request.Source.Password)
		ref = getRef(destination, request)
	}

	versions := []models.Version{}
	versions = append(versions, models.Version{
		Sha: ref,
	})

	json.NewEncoder(os.Stdout).Encode(versions)
}

//getRef ...
//returns the ref of the ready branch
func getRef(basePath string, request models.CheckRequest) (ref string) {

	chErr := os.Chdir(basePath)
	if chErr != nil {
		log.Panicln("no basepath", basePath, chErr)
	}

	if err := githandler.HardReset(); err != nil {
		log.Panicln(err)
	}

	if err := githandler.Pull(); err != nil {
		log.Panicln("could not pull from remote: ", err)
	}

	if err := githandler.FetchPrune(); err != nil {
		log.Panicln("could not fetch from remote: ", err)
	}

	branchName := phlow.UpNext("origin", request.Source.PrefixReady)
	if branchName == "" {
		fmt.Fprintln(os.Stderr, "No ready branches found")

		if request.Version.Sha != "" {
			return request.Version.Sha
		}

		//First build with no ready branches
		fmt.Fprintln(os.Stderr, "Create ready branch for this error to go away")
		os.Exit(1)

	}

	err := githandler.CheckOut(branchName)
	if err != nil {
		log.Panicln(err)
	}

	if ref, err = githandler.RevParse(); err != nil {
		log.Panicln(err)
	}

	err = githandler.CheckOut(request.Source.MainBranch)
	if err != nil {
		log.Panicln(err)

	}
	return ref
}

//doesExist ...
//checks the repository is still cached
func doesExist(basePath string) bool {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return false
	}
	return true
}
