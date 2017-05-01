package repo

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/praqma/concourse-git-phlow/githandler"
)

//WriteRDYBranch ...
//writes the name of the branch to the file
func WriteRDYBranch(name string) {
	err := ioutil.WriteFile(".git/git-phlow-ready-branch", []byte(name), 0655)
	if err != nil {
		fmt.Fprintln(os.Stderr,"Could not file for ready branch")
		os.Exit(1)
	}
}

//RenameRemoteBranch ...
//renames a branch with a prefix
func RenameRemoteBranch(URL, newName, oldName string) {
	err := githandler.PushRenameHTTPS(URL, newName, oldName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not push rename old: %s new: %s", oldName, newName)
		os.Exit(1)
	}

	err = githandler.PushDeleteHTTPS("origin", oldName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not push delete branch %s \n", oldName)
		os.Exit(1)
	}
}

//FormatURL ...
//formats the https url
func FormatURL(URL, username, password string) string {
	ct := strings.Replace(URL, "https://", "", 1)
	u := fmt.Sprintf("https://%s:%s@%s", username, password, ct)
	return u
}

//CloneRepoSource ...
//clones the repository
func CloneRepoSource(URL, path, username, password string) {
	c := FormatURL(URL, username, password)
	fmt.Fprintf(os.Stderr, "Cloning into desitnation: %s from:  %s\n", path, URL)
	_, err := githandler.Clone(c, path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get repo from: %s", URL)
		os.Exit(1)
	}
}

func Check(e error, str string) {
	if e != nil {
		fmt.Fprintln(os.Stderr, str)
		os.Exit(1)
	}
}
