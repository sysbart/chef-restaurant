package chef-restaurant

import (
 "strings"
 "flag"
 "./chef/"
 "./git/"
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
			chef.Upload("cookbook", file)
		}
	}
}
