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
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <source>")
		os.Exit(1)
	}

	destination := os.Args[1]

	var request models.OutRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "OS in parsing errored")
		os.Exit(1)
	}

	err = os.Chdir(destination + "/" + request.Params.Repository)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not change dir:", err.Error())
		os.Exit(1)
	}

	name, err := ioutil.ReadFile(".git/git-phlow-ready-branch")
	repo.Check(err, "failed reading branch name from file")

	fmt.Fprintln(os.Stderr, string(name))

	if string(name) == "" {
		fmt.Fprintln(os.Stderr, "No ready branch to integrate with master.. Exiting build")
		SendMetadata()
		os.Exit(0)
	}

	HttpsPush(request.Source.URL, request.Source.Username, request.Source.Password)

	err = githandler.PushDeleteHTTPS("origin", string(name))
	if err != nil {
		fmt.Fprintln(os.Stderr, "branch could not be deleted:", err.Error())
		os.Exit(1)
	}

	SendMetadata()
}

//SendMetadata ...
func SendMetadata() {
	ref, _ := githandler.CommitSha()
	author, _ := githandler.Author()
	date, _ := githandler.AuthorDate()

	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version: models.Version{Sha: ref},
		MetaData: models.Metadata{
			{"commit", ref},
			{"author", author},
			{"authordate", date},
		},
	})
}

func HttpsPush(URL string, username, password string) {

	url := repo.FormatURL(URL, username, password)

	fmt.Fprintf(os.Stderr, "pushing to: %s \n", URL)
	_, err := githandler.PushHTTPS(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not push to repository")
		os.Exit(1)
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
