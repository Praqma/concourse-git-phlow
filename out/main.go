package main

import (
	"fmt"
	"os"
	"github.com/groenborg/pip/models"
	"encoding/json"
	"github.com/groenborg/pip/githandler"
	"strings"
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

	githandler.LS()

	out, _ := githandler.Branch()
	fmt.Fprintln(os.Stderr, out)
	ref, _ := githandler.RevParse()

	HttpsPush(request.Source.URL, request.Source.Username, request.Source.Password)

	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version: models.Version{Sha: ref},
		MetaData: models.Metadata{
			{"author", "david"},
		},
	})
}

func HttpsPush(URL string, username, password string) {
	ct := strings.Replace(URL, "https://", "", 1)
	pu := fmt.Sprintf("https://%s:%s@%s", username, password, ct)

	fmt.Fprintf(os.Stderr, "pushing to: %s \n", URL)
	_, err := githandler.PushHTTPS(pu)
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
