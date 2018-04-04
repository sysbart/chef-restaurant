package main

import (
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

func processFilePerObjectType(file string) (object string, err error){
	if strings.Contains(file, "environments/") {
		object = "environment"
		err = chef.Upload("environment", file)
	} else if strings.Contains(file, "roles/") {
		object = "role"
		err = chef.Upload("role", file)
	} else if strings.Contains(file, "cookbooks/") && !strings.Contains(file, "fixtures/") {
		object = "cookbook"
		bumpLevel := "patch"
		cookbookNewVersion := helpers.Metadata(file, bumpLevel)
		cookbookName, _, err := chef.CookbookInfo(file)
		if err != nil {
			return object, err
		}

		err = git.Commit(file, fmt.Sprintf("Pin %s version for the cookbook %s", cookbookNewVersion, cookbookName))
		if err != nil {
			return object, err
		}

		err = git.Push()
		if err != nil{
			return object, err
		}

		err = chef.Upload("cookbook", file)
		if err != nil{
			return object, err
		}
	} else {
		object = ""
	}

	return object, err
}

func main() {
	var err error

	options := config.GenerateOptions()
	config := config.Init(options.ConfigFolder)
	log.SetLevel(options.LogLevelMethod)

	helpers.WorkingFolder = options.WorkingFolder
	notification.SlackNotificationHookURL = config.SlackNotificationHookURL
	notification.SlackNotificationChannel = config.SlackNotificationChannel
	git.GitHubOrganizationName = config.GitHubOrganizationName
	git.GitHubRepoName = config.GitHubRepoName

	if options.Commit == "" {
		options.Commit, err = git.LastCommit()
		if err != nil {
			log.Fatal("Unable to get last commit id, exiting.")
		}
	}

	commitInfo, err := git.CommitInfo(options.Commit)
	if err != nil {
		log.Fatal("Unable to get last commit information, exiting.")
	}

	if strings.Contains(commitInfo.Author, config.SkippedAuthorName) {
		log.Info("The author of the last commit is chef-restaurant. I am exiting.")
		os.Exit(0)
	}
	// Environment and roles are file based, cookbooks are folder based
	files, err := git.FilesListForEachCommit(options.Commit)
	if err != nil {
		log.Fatal("Unable to get the list of files for all commits, exiting.")
	}

	for _, file := range files {
			var err error
			var object string
			var notificationTitle string
			var notificationMessage string
			var notificationColor string

			_, err = os.Stat(file)
			if err == nil {
				log.Debugf("Processing file: %s", file)
				object, err = processFilePerObjectType(file)
			} else {
				log.Debugf("The file %s was skipped since it was deleted on Git", file)
				continue
			}

			if object != "" {
				if err == nil {
					if object == "cookbook" {
						_, cookbookVersion, _ := chef.CookbookInfo(file)
						notificationTitle = fmt.Sprintf("%s *%s* [%s] has been uploaded to the Chef server", strings.Title(object), chef.ParseObjectByFileName(object, file), cookbookVersion)
					} else {
						notificationTitle = fmt.Sprintf("%s *%s* has been uploaded to the Chef server", strings.Title(object), chef.ParseObjectByFileName(object, file))
					}
					notificationMessage = fmt.Sprintf("`<%s|%s>` %s - [%s] ", git.GenerateCommitURL(options.Commit), options.Commit, commitInfo.Title, commitInfo.Author)
					notificationColor = "good"
				} else {
					notificationTitle = fmt.Sprintf("[ERROR] %s *%s* has been not uploaded to the Chef server", strings.Title(object), chef.ParseObjectByFileName(object, file))
					notificationMessage = fmt.Sprintf("`<%s|%s>` %s - [%s] ", git.GenerateCommitURL(options.Commit), options.Commit, commitInfo.Title, commitInfo.Author)
					notificationColor = "danger"
				}

				log.Infof("Sending notification for file: %s", file)
				log.Debugf("Notification metadata: %s, %s, %s", notificationTitle, notificationMessage, notificationColor)

				notification.SendMessage(notificationTitle, notificationMessage, notificationColor)
			} else {
			log.Debugf("Skipping notification for file: %s", file)
		}
	}
}
