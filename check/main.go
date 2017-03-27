package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/groenborg/pip/models"
	"github.com/groenborg/pip/githandler"
)

func main() {

	var request models.CheckRequest
	var ref string
	basePath := os.Getenv("TMPDIR") + "/cache"

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}


	if doesExist(basePath) {
		ref = GetRef(basePath)
	} else {
		getRepo(basePath, request.Source.URL)
	}


	versions := []models.Version{}
	versions = append(versions, models.Version{Sha: ref})

	json.NewEncoder(os.Stdout).Encode(versions)

}

func GetRef(basePath string) (ref string) {
	os.Chdir(basePath)
	if err := githandler.Fetch(); err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err.Error())
		os.Exit(1)
	}

	var err error
	if ref, err = githandler.RevParse(); err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err.Error())
		os.Exit(1)
	}
	return ref

}

func getRepo(basePath, url string) {

	_, err := githandler.Clone(url, basePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err.Error())
		os.Exit(1)
	}
}

func doesExist(basePath string) bool {

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return false
	}
	return true
}
