package main

import (
	"fmt"
	"os"
	"github.com/groenborg/pip/models"
	"encoding/json"
	"github.com/groenborg/pip/githandler"
	"path/filepath"
	"io"
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

	file, err := os.Create(filepath.Join(destination, "cloned-repo"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "create file fails:", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	var request models.InRequest
	var sha string

	err = json.NewDecoder(io.TeeReader(os.Stdin, file)).Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "OS in parsing errored")
		os.Exit(1)
	}

	sha = request.Version.Sha

	err = os.Chdir(destination)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not change dir:", err.Error())
		os.Exit(1)
	}

	getRepo(request.Source.URL)

	os.Chdir("./phlow-test")

	err = githandler.Merge(sha)
	if err != nil {
		fmt.Fprintln(os.Stderr, "checkout fail:", err.Error())
		os.Exit(1)
	}

	GetMeteData()
}

func GetMeteData() {
	sha, err := githandler.CommitSha()
	if err != nil {
		fmt.Fprintln(os.Stderr, "commit sha failed:", err.Error())
		os.Exit(1)
	}

	author, err := githandler.Author()
	if err != nil {
		fmt.Fprintln(os.Stderr, "author failed:", err.Error())
		os.Exit(1)
	}

	date, err := githandler.AuthorDate()
	if err != nil {
		fmt.Fprintln(os.Stderr, "date failed:", err.Error())
		os.Exit(1)
	}

	metadata := []models.MetaData{}
	metadata = append(metadata, models.MetaData{Author: author, Commit: sha, AuthorDate: date})

	json.NewEncoder(os.Stdout).Encode(metadata)
}

func getRepo(url string) {

	_, err := githandler.CloneCurrentDir(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "get repo failed:", err.Error())
		os.Exit(1)
	}
}
