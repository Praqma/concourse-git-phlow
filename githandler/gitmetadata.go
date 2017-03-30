package githandler

import (
	"strings"
	"github.com/praqma/concourse-git-phlow/executor"
)

func Author() (author string, err error) {
	return executor.ExecuteCommand("git", "log", "-1", "--format=format:%an")
}

func CommitSha() (sha string, err error) {
	sha, err = executor.ExecuteCommand("git", "rev-parse", "HEAD")
	sha = strings.TrimSpace(sha)
	return

}
func AuthorDate() (date string, err error) {
	return executor.ExecuteCommand("git", "log", "-1", "--format=format:%ai")
}
