package main

import (
	"fmt"
	"os"
	"github.com/groenborg/pip/models"
	"encoding/json"
	"github.com/groenborg/pip/githandler"
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

	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version: models.Version{Sha: request.Version.Sha},
		MetaData: models.Metadata{
			{"author", "david"},
		},
	})
}
