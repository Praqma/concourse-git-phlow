package githandler

import (
	"github.com/groenborg/pip/executor"
	"strings"
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
