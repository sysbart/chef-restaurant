package main

import (
	"flag"
	"fmt"
	"github.com/sysbart/chef-restaurant/chef"
	"github.com/sysbart/chef-restaurant/config"
	"github.com/sysbart/chef-restaurant/git"
	"github.com/sysbart/chef-restaurant/helpers"
	"github.com/sysbart/chef-restaurant/notification"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	optCommit := flag.String("commit", "", "commit ID")
	optConfigFolder := flag.String("configFolder", "", "config directory")
	flag.Parse()

	commit := *optCommit
	if commit == "" {
		commit = git.LastCommit()
	}

	configFolder := *optConfigFolder
	if configFolder == "" {
		configFolder = "etc"
	}

	commitAuthor := git.CommitInfo(commit).Author
	commitTitle := git.CommitInfo(commit).Title

	config := config.Init(configFolder)
	notification.SlackNotificationHookURL = config.SlackNotificationHookURL
	notification.SlackNotificationChannel = config.SlackNotificationChannel
	git.GitHubOrganizationName = config.GitHubOrganizationName
	git.GitHubRepoName = config.GitHubRepoName

	if strings.Contains(commitAuthor, config.SkippedAuthorName) {
		log.Info("The author of the last commit is chef-restaurant. I am exiting.")
		os.Exit(0)
	}
	// Environment and roles are file based, cookbooks are folder based
	files := git.FilesListForEachCommit(commit)
	for _, file := range files {
		var object string
		var notify bool
		if strings.Contains(file, "environments/") {
			notify = true
			object = "environment"
			chef.Upload("environment", file)
		} else if strings.Contains(file, "roles/") {
			notify = true
			object = "role"
			chef.Upload("role", file)
		} else if strings.Contains(file, "cookbooks/") {
			notify = true
			object = "cookbook"
			bumpLevel := "patch"
			newVersion := helpers.Metadata(file, bumpLevel)
			cookbookName := chef.CookbookInfo(file)
			git.Commit(file, fmt.Sprintf("Pin %s version for the cookbook %s", newVersion, cookbookName))
			git.Push()
			chef.Upload("cookbook", file)
		}

		if notify {
			notificationTitle := fmt.Sprintf("%s *%s* has been uploaded to the Chef server", strings.Title(object), chef.ParseObjectByFileName(object, file))
			notificationMessage := fmt.Sprintf("`<%s|%s>` %s - [%s] ", git.GenerateCommitURL(commit), commit, commitTitle, commitAuthor)
			notification.SendMessage(notificationTitle, notificationMessage)
		}

	}
}
