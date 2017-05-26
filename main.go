package main

import (
 "strings"
 "flag"
 "fmt"
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

	files := git.FilesListForEachCommit(commit)
	git.CommitInfo(commit)

	// Environment and roles are file based, cookbooks are folder based
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
			chef.Upload(cookbookName, file)
		}
	}
}
