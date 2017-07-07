package git

import (
	"fmt"
	"github.com/sysbart/chef-restaurant/helpers"
	"github.com/sysbart/chef-restaurant/structs"
	"regexp"
	"strings"
)

var GitHubOrganizationName string
var GitHubRepoName string

func Push() {
	helpers.RunCommand("git", "push", "origin", "master")
}

func Commit(file string, message string) {
	helpers.RunCommand("git", "commit", file, "-m", message)

	fmt.Printf("Committed the file %s [%s]\n", file, message)
}

func LastCommit() string {
	commitCmd := helpers.RunCommand("git", "log", "-n1", "--pretty=format:%h")

	commit := string(commitCmd)
	fmt.Printf("Last commit found : %s\n", commit)

	return commit
}

func CommitInfo(commit string) structs.CommitInfo {
	merge := false
	commitRegexp := regexp.MustCompile(`(?s)commit (.*)Author: (.*)Date: (.*)\n\n(.*)`)

	commitCmd := helpers.RunCommand("git", "log", "-n1", commit)

	if strings.Contains(string(commitCmd), "^Merge:") {
		commitRegexp = regexp.MustCompile(`(?s)commit (.*)Merge: .*Author: (.*)Date: (.*)\n\n(.*)`)
		merge = true
	}

	commitSubMatch := commitRegexp.FindStringSubmatch(string(commitCmd))
	commit, author, date, title := strings.TrimSpace(commitSubMatch[1]), strings.TrimSpace(commitSubMatch[2]), strings.TrimSpace(commitSubMatch[3]), strings.TrimSpace(commitSubMatch[4])

	return structs.CommitInfo{merge, author, date, title}
}

func GenerateCommitURL(commit string) string {
	URL := fmt.Sprintf("https://github.com/%s/%s/commit/%s", GitHubOrganizationName, GitHubRepoName, commit)
	return URL
}

func FilesListForEachCommit(commit string) []string {
	fmt.Printf("Retrieving files list for the commit %s\n", commit)
	filesCmd := helpers.RunCommand("git", "diff", "--name-only", commit, commit+"^")
	filesCmdUnFiltered := string(filesCmd)
	cookbookRegexp := regexp.MustCompile(`(.*cookbooks/[^/]*)/(.*)`)
	filesCmdFiltered := cookbookRegexp.ReplaceAllString(filesCmdUnFiltered, "$1")
	filesList := helpers.RemoveDuplicatesUnordered(strings.Split(strings.TrimSpace(filesCmdFiltered), "\n"))

	return filesList
}
