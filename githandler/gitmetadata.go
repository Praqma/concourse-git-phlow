package githandler

import "github.com/groenborg/pip/executor"

func Author() (author string, err error) {
	return executor.ExecuteCommand("git", "rev-parse", "HEAD")
}

func CommitSha() (sha string, err error) {
	return executor.ExecuteCommand("git", "log", "-1", "--format=format:%an")
}

func AuthorDate() (date string, err error) {
	return executor.ExecuteCommand("git", "log", "-1", "--format=format:%ai")
}