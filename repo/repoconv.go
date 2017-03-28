package repo

import (
	"strings"
	"fmt"
	"os"
	"github.com/groenborg/pip/githandler"
)

func FormatURL(URL, username, password string) string {
	ct := strings.Replace(URL, "https://", "", 1)
	u := fmt.Sprintf("https://%s:%s@%s", username, password, ct)
	return u
}

func CloneRepoSource(URL, path, username, password string) {
	c := FormatURL(URL, username, password)
	fmt.Fprintf(os.Stderr, "Cloning into desitnation: %s from:  %s\n", path, URL)
	_, err := githandler.Clone(c, path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "get repo failed")
		os.Exit(1)
	}
}
