package githandler

import (
	"regexp"
	"strings"

	"fmt"
	"os"

	"github.com/praqma/git-phlow/executor"
)

//RebaseOnto ...
//Rebase current branch onto delivery branch
func RebaseOnto(br string) (err error) {
	_, err = executor.ExecuteCommand("git", "rebase", br)
	return
}

//MergeFFO ...
//Only merge if it is a fast forward
func MergeFFO(branch string) (err error) {
	_, err = executor.ExecuteCommand("git", "merge", "--ff-only", branch)
	return
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

	str, err := executor.ExecuteCommand("git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	str = strings.TrimSpace(str)
	return str, nil
}

//Fetch ...
func Status() string {
	out, _ := executor.ExecuteCommand("git", "branch", "-av")
	return out
}

//Fetch ...
func Fetch() error {
	_, err := executor.ExecuteCommand("git", "fetch", "--prune")
	return err
}

func PushRenameHTTPS(URL string, new, old string) (err error) {
	rn := fmt.Sprintf("%s:%s", old, new)
	fmt.Fprintln(os.Stderr, rn)
	_, err = executor.ExecuteCommand("git", "push", "--repo", URL, "origin", strings.TrimSpace(rn))
	return
}

func PushDeleteHTTPS(remote, name string) (err error) {
	_, err = executor.ExecuteCommand("git", "push", remote, "--delete", name)
	return
}

//PushHTTPS ...
func PushHTTPS(URL string) (string, error) {
	return executor.ExecuteCommand("git", "push", "--repo", URL)
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
