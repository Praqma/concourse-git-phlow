package main

import (
	"fmt"
	"os"
	"github.com/groenborg/pip/models"
	"encoding/json"
	"github.com/groenborg/pip/githandler"
	"bufio"
)

func main() {

	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

	destination := os.Args[1]

	err := os.MkdirAll(destination, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mkdir all fails:", err.Error())
		os.Exit(1)
	}

	var request models.InRequest

	fmt.Println(os.Stderr, os.Stdin)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	fmt.Fprintln(os.Stderr, text)

	err = json.Unmarshal([]byte(text), &request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "OS in parsing errored")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, request.Source.URL)
	fmt.Fprintln(os.Stderr, request.Version.Sha)

	getRepo(request.Source.URL, destination)

	err = os.Chdir(destination)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not change dir:", err.Error())
		os.Exit(1)
	}

	GetMeteData(request.Version.Sha)
}

func GetMeteData(sha string) {
	sha, err := githandler.CommitSha()
	if err != nil {
		fmt.Fprintln(os.Stderr, "commit sha failed:", err.Error())
		os.Exit(1)
	}

	//author, err := githandler.Author()
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "author failed:", err.Error())
	//	os.Exit(1)
	//}
	//
	//date, err := githandler.AuthorDate()
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "date failed:", err.Error())
	//	os.Exit(1)
	//}

	githandler.LS()

	fmt.Fprintln(os.Stderr, "ABOUT TO RETURN JSON")

	metadata := []models.MetaData{}
	metadata = append(metadata, models.MetaData{Name: "commit", Value: sha})


	req := models.OutRequest{Version: models.Version{Sha: sha}, MetaData: metadata}

	fmt.Fprintln(os.Stderr, req)

	json.NewEncoder(os.Stdout).Encode(req)
}

func getRepo(url, path string) {

	_, err := githandler.Clone(url, path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "get repo failed:", err.Error())
		os.Exit(1)
	}
}
