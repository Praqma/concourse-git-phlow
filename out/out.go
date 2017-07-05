package main

import (
	"fmt"
	"os"

	"encoding/json"

	"strings"

	"io/ioutil"
	"github.com/praqma/concourse-git-phlow/githandler"
	"github.com/praqma/concourse-git-phlow/models"
	"github.com/praqma/concourse-git-phlow/repo"
	"github.com/praqma/concourse-git-phlow/concourse"
	"log"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <source>")
		os.Exit(1)
	}

	destination := os.Args[1]

	var request models.OutRequest

	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Panicln(err)
	}

	fmt.Fprintln(os.Stderr, "Resource Version "+repo.Version)

	if err := os.Chdir(destination + "/" + request.Params.Repository); err != nil {
		log.Panicln(err, destination)
	}

	name, err := ioutil.ReadFile(".git/git-phlow-ready-branch")
	if err != nil {
		log.Panicln(err, name)
	}

	if string(name) == "" || !BranchExistsOnOrigin(string(name)) {
		fmt.Fprintln(os.Stderr, "No ready branch to integrate with master.. Exiting build")
		fmt.Fprintln(os.Stderr, "Output")
		ref, _ := githandler.CommitSha()
		concourse.SendMetadata(ref)
		os.Exit(0)
	}

	HttpsPush(request.Source.URL, request.Source.Username, request.Source.Password)

	err = githandler.PushDeleteHTTPS("origin", string(name))
	if err != nil {
		log.Panicln(err, "branch could not be deleted")

	}
	fmt.Fprintf(os.Stderr, "%s has been pushed to %s", string(name), request.Source.MainBranch)
	ref, _ := githandler.CommitSha()
	concourse.SendMetadata(ref)
}

//BranchExistsOnOrigin ...
func BranchExistsOnOrigin(branchName string) (exists bool) {
	branchName = strings.TrimSpace(branchName)

	if err := githandler.FetchPrune(); err != nil {
		log.Panicln(err)
	}

	brOut, err := githandler.BranchList()
	if err != nil {
		log.Panicln(err)
	}

	var list []string
	for _, branch := range strings.Split(brOut, "\n") {
		if branch != "" {
			branch = strings.TrimSpace(branch)
			list = append(list, branch)
		}
	}

	for _, branch := range list {
		if strings.Contains(branch, branchName) {
			return true
		}
	}
	return false
}

func HttpsPush(URL string, username, password string) {

	url := repo.FormatURL(URL, username, password)

	fmt.Fprintf(os.Stderr, "pushing to: %s \n", URL)
	_, err := githandler.PushHTTPS(url)
	if err != nil {
		log.Panicln(err, URL)
	}
}

func remoteURLExtractor(url string) (ssh bool, http bool) {
	//Extracts repo and org from ssh url format
	if strings.HasPrefix(url, "git@") {
		return true, false
	}
	//Extracts repo and org from http url format
	if strings.HasPrefix(url, "https") {
		return false, true
	}
	//Clone from local repo
	return false, false
}
