package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/groenborg/pip/models"
	"github.com/groenborg/pip/githandler"
	"github.com/groenborg/pip/repo"
)

func main() {
	var request models.CheckRequest
	var ref string
	destination := os.Getenv("TMPDIR") + "/cache"

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not parse input in check")
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
	os.Chdir(basePath)


	if err := githandler.Fetch(); err != nil {
		fmt.Fprintln(os.Stderr, "could not fetch from remote: ", err.Error())
		os.Exit(1)
	}

	branchName, err := githandler.PhlowReadyBranch()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed getting ready branch: ", err.Error())
		os.Exit(1)
	}

	err = githandler.CheckOut(branchName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "checkout failed: ", err.Error())
		os.Exit(1)
	}

	if ref, err = githandler.RevParse(); err != nil {
		fmt.Fprintln(os.Stderr, "could not retrieve ref:: ", err.Error())
		os.Exit(1)
	}

	err = githandler.CheckOut(request.Source.Master)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not checkout main branch:", err.Error())
		os.Exit(1)
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
