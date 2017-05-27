package main

import (
 "strings"
 "flag"
 "fmt"
 "os"
 "./chef/"
 "./git/"
 "./helpers/"
)

func main() {
	optCommit := flag.String("commit", "", "commit ID")
	flag.Parse()

	commit := *optCommit
	if commit == "" {
	        commit = git.LastCommit()
	}

	commitAuthor := git.CommitInfo(commit).Author
  if strings.Contains(commitAuthor, "chef-restaurant") {
    fmt.Println("The author of the last commit is chef-restaurant. I am exiting.")
    os.Exit(0)
  }

	// Environment and roles are file based, cookbooks are folder based
  files := git.FilesListForEachCommit(commit)
	for _, file := range files {
		if strings.Contains(file, "environments/") {
			chef.Upload("environment", file)
		} else if strings.Contains(file, "roles/") {
			chef.Upload("role", file)
		} else if strings.Contains(file, "cookbooks/") {

      bumpLevel := "patch"
      newVersion := helpers.Metadata(file, bumpLevel)
      cookbookName := chef.CookbookInfo(file)

      git.Commit(file, fmt.Sprintf("Pin %s version for the cookbook %s", newVersion, cookbookName))
      git.Push()
			chef.Upload("cookbook", file)
		}
	}
}
