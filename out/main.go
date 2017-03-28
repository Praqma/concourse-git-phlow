package main

import (
	"fmt"
	"os"
	"github.com/groenborg/pip/models"
	"encoding/json"
	"github.com/groenborg/pip/githandler"
	"strings"
	"github.com/groenborg/pip/auth"
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

	fmt.Fprintln(os.Stderr, request)
	fmt.Fprintln(os.Stderr, destination)

	err = os.Chdir(destination + "/" + request.Params.Repository)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not change dir:", err.Error())
		os.Exit(1)
	}

	name, err := githandler.PhlowReadyBranch()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Getting ready branch fail:", err.Error())
		os.Exit(1)
	}

	HttpsPush(request.Source.URL, request.Source.Username, request.Source.Password)

	_, err = githandler.BranchDelete(name, "origin")
	if err != nil {
		fmt.Fprintln(os.Stderr, "branch could not be deleted:", err.Error())
	}

	ref, _ := githandler.RevParse()
	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version: models.Version{Sha: ref},
		MetaData: models.Metadata{
			{"author", "david"},
		},
	})
}

func HttpsPush(URL string, username, password string) {

	url := auth.FormatURL(URL,username,password)

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
