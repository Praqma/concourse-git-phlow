package mwriter

import (
	"encoding/json"
	"github.com/praqma/concourse-git-phlow/models"
	"fmt"
	"os"
	"github.com/praqma/concourse-git-phlow/githandler"
)

//GetMetadata ...
//sends the metadata to output
func SendMetadata(sha string) {
	ref, _ := githandler.CommitSha()
	author, _ := githandler.Author()
	date, _ := githandler.AuthorDate()

	str, err := json.Marshal(models.InResponse{
		Version: models.Version{Sha: sha},
		MetaData: models.Metadata{
			{"commit", ref},
			{"author", author},
			{"authordate", date},
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Fprint(os.Stdout, string(str))
	//json.NewEncoder(os.Stdout).Encode()
}
