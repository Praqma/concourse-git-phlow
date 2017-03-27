package githandler

import (
	"regexp"
	"strings"

	"github.com/groenborg/git-phlow/executor"
)

func PhlowReadyBranch() (name string, err error) {
	return executor.ExecuteCommand("git", "phlow", "agent", "upnext")
}

//Clone ...
func Clone(URL string, path string) (output string, err error) {
	output, err = executor.ExecuteCommand("git", "clone", URL, path)
	return
}

//CheckOut ...
func CheckOut(branch string) error {
	_, err := executor.ExecuteCommand("git", "checkout", branch)
	return err
}

func RevParse() (out string, err error) {
	return executor.ExecuteCommand("git", "rev-parse", "HEAD")
}

//Status ...
func Status() error {
	_, err := executor.ExecuteCommand("git", "status")
	return err
}

//Add ...
func Add() error {
	_, err := executor.ExecuteCommand("git", "add", "--all")
	return err
}

//Commit ...
func Commit(message string) (string, error) {
	return executor.ExecuteCommand("git", "commit", "-m", message)
}

//Fetch ...
func Fetch() error {
	_, err := executor.ExecuteCommand("git", "fetch", "--all")
	return err
}

//FetchPrune ...
func FetchPrune() error {
	_, err := executor.ExecuteCommand("git", "fetch", "--prune")
	return err
}

//Pull ...
func Pull() (string, error) {
	return executor.ExecuteCommand("git", "pull", "--rebase")
}

//Push ...
func Push() (string, error) {
	return executor.ExecuteCommand("git", "push")
}

//Merge ...
func Merge(branch string) error {
	_, err := executor.ExecuteCommand("git", "merge", branch)
	return err
}

//RemoteInfo ...
type RemoteInfo struct {
	Organisation string
	Repository   string
}

//Remote ...
//Must have either origin or upstream
//THIS NEEDS TO BE REVISITED
func Remote() (*RemoteInfo, error) {
	var res string
	var err error
	if res, err = executor.ExecuteCommand("git", "ls-remote", "--get-url", "origin"); err != nil {
		return nil, err
	}
	res = strings.Trim(res, "\n")
	return remoteURLExtractor(res), nil
}

func remoteURLExtractor(url string) *RemoteInfo {
	re := regexp.MustCompile(`.+:(\S+)\/(\S+)\.git`)

	//Extracts repo and org from ssh url format
	if strings.HasPrefix(url, "git@") {
		match := re.FindStringSubmatch(url)
		return &RemoteInfo{match[1], match[2]}
	}
	//Extracts repo and org from http url format
	if strings.HasPrefix(url, "http") {
		splitURL := strings.Split(strings.TrimSuffix(url, ".git"), "/")
		org := splitURL[len(splitURL)-2]
		repo := splitURL[len(splitURL)-1]
		return &RemoteInfo{org, repo}
	}

	//Clone from local repo
	return &RemoteInfo{}
}
