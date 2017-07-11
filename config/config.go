package config

import (
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

type Config struct {
  SkippedAuthorName         string
  SlackNotificationHookURL  string
  SlackNotificationChannel  string
  GitHubOrganizationName    string
  GitHubRepoName            string
}

func Init(configFolder string) Config {
	viper.AddConfigPath(configFolder)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	skippedAuthorName := viper.GetString("skipped_author_name")
	slackNotificationHookURL := viper.GetString("slack.notification.hook_url")
	slackNotificationChannel := viper.GetString("slack.notification.channel")
	gitHubOrganizationName := viper.GetString("github.organization")
  gitHubRepoName := viper.GetString("github.repo")

	if slackNotificationHookURL == "" || slackNotificationChannel == "" {
		log.Fatal("Slack configuration notification is not setup. I am exiting.")
	}

	if gitHubOrganizationName == "" || gitHubRepoName == "" {
		log.Fatal("Github configuration is not setup. I am exiting.")
  }

  return Config{skippedAuthorName,slackNotificationHookURL,slackNotificationChannel,gitHubOrganizationName,gitHubRepoName}
}
