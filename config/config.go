package config

import (
	"flag"
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

type Options struct {
	Commit         		string
	ConfigFolder  		string
	WorkingFolder  		string
	LogLevelMethod  	log.Level
}

func GenerateOptions() Options {
	optCommit := flag.String("commit", "", "commit ID")
	optConfigFolder := flag.String("configFolder", "etc", "config directory")
	optWorkingFolder := flag.String("workingFolder", "", "working directory")
	optLogLevel := flag.String("logLevel", "info", "log level")

	flag.Parse()

	var logLevelMethod log.Level
	switch *optLogLevel {
	case "debug":
		logLevelMethod = log.DebugLevel
	case "warn":
		logLevelMethod = log.WarnLevel
	case "info":
		logLevelMethod = log.InfoLevel
	}

	return Options{*optCommit, *optConfigFolder, *optWorkingFolder, logLevelMethod}
}
