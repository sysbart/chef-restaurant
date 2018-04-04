package git

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sysbart/chef-restaurant/helpers"
	"github.com/sysbart/chef-restaurant/structs"
	"regexp"
	"strings"
)

var GitHubOrganizationName string
var GitHubRepoName string

func Push() error {
	_, err := helpers.RunCommand("git", "push", "origin", "master")

	if err != nil {
		log.Errorf("Git failed to push this commit to the branch")
		return err
	}

	log.Info("Pushed the commit to the branch")
	return nil
}

func Commit(file string, message string) error {
	_, err := helpers.RunCommand("git", "commit", file, "-m", message)

	if err != nil {
		log.Errorf("Git failed to push this commit to the branch")
		return err
	}

	log.Infof("Committed the file %s [%s]\n", file, message)
	return nil
}

func LastCommit() (string, error) {
	var commit string

	commitCmd, err := helpers.RunCommand("git", "log", "-n1", "--pretty=format:%h")

	if err != nil {
		log.Errorf("Git failed to get the last commit of this branch")
		return commit, err
	}

	commit = string(commitCmd)
	log.Infof("Last commit found : %s\n", commit)
	return commit, nil
}

func CommitInfo(commit string) (structs.CommitInfo, error) {
	merge := false
	var commitInfo structs.CommitInfo

	commitRegexp := regexp.MustCompile(`(?s)commit (.*)Author: (.*)Date: (.*)\n\n(.*)`)

	commitCmd, err := helpers.RunCommand("git", "log", "-n1", commit)
	if err != nil {
		log.Errorf("Git failed to get details about the last commit of this branch")
		return commitInfo, err
	}

	if strings.Contains(string(commitCmd), "^Merge:") {
		commitRegexp = regexp.MustCompile(`(?s)commit (.*)Merge: .*Author: (.*)Date: (.*)\n\n(.*)`)
		merge = true
	}

	commitSubMatch := commitRegexp.FindStringSubmatch(string(commitCmd))

	commitInfo.IsMergeCommit = merge
	commitInfo.Author = strings.TrimSpace(commitSubMatch[2])
	commitInfo.Date = strings.TrimSpace(commitSubMatch[3])
	commitInfo.Title = strings.TrimSpace(commitSubMatch[4])

	return commitInfo, nil
}

func GenerateCommitURL(commit string) string {
	URL := fmt.Sprintf("https://github.com/%s/%s/commit/%s", GitHubOrganizationName, GitHubRepoName, commit)
	return URL
}

func FilesListForEachCommit(commit string) ([]string, error) {
	var filesList []string
	var err error

	log.Infof("Retrieving files list for the commit %s", commit)
	filesCmd, err := helpers.RunCommand("git", "diff", "--name-only", commit, commit+"^")
	if err != nil {
		log.Errorf("Git failed to get the list of files for commit %s", commit)
		return filesList, err
	}

	filesCmdUnFiltered := string(filesCmd)
	cookbookRegexp := regexp.MustCompile(`(.*cookbooks/[^/]*)/(.*)`)
	filesCmdFiltered := cookbookRegexp.ReplaceAllString(filesCmdUnFiltered, "$1")
	filesList = helpers.RemoveDuplicatesUnordered(strings.Split(strings.TrimSpace(filesCmdFiltered), "\n"))

	return filesList, nil
}
