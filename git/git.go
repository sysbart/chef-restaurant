package git

import (
 "os/exec"
 "fmt"
 "log"
 "strings"
 "regexp"
 "../helpers"
 "../structs"
)


func Push() {
	_, err := exec.Command("git", "push", "origin", "master").Output()
	if err != nil {
		log.Fatal(err)
	}
}


func Commit(file string, message string) {
	_, err := exec.Command("git", "commit", file, "-m", message).Output()
	if err != nil {
		log.Fatal(err)
	}

  fmt.Printf("Committed the file %s [%s]", file, message)
}

func LastCommit() string {
	commitCmd, err := exec.Command("git", "log", "-n1", "--pretty=format:%h").Output()
	if err != nil {
		log.Fatal(err)
	}

	commit := string(commitCmd)
	fmt.Println("Last commit found: " + commit)
	return commit
}

func CommitInfo(commit string) structs.CommitInfo {
  merge := false
  commitRegexp := regexp.MustCompile(`(?s)commit (.*)Author: (.*)Date: (.*)\n\n(.*)`)

	commitCmd, err := exec.Command("git", "log", "-n1").Output()
	if err != nil {
		log.Fatal(err)
	}

  if strings.Contains(string(commitCmd), "^Merge:") {
    commitRegexp = regexp.MustCompile(`(?s)commit (.*)Merge: .*Author: (.*)Date: (.*)\n\n(.*)`)
    merge = true
  }

  commitSubMatch := commitRegexp.FindStringSubmatch(string(commitCmd))
  commit, author, date, title := strings.TrimSpace(commitSubMatch[1]), strings.TrimSpace(commitSubMatch[2]), strings.TrimSpace(commitSubMatch[3]), strings.TrimSpace(commitSubMatch[4])

	return structs.CommitInfo{merge, author, date, title}
}

func FilesListForEachCommit(commit string) []string {
	fmt.Println("Retrieving files list for the commit " + commit )
	filesCmd, err := exec.Command("git", "diff", "--name-only", commit, commit + "^").Output()
	if err != nil {
		log.Fatal(err)
	}
	filesCmdUnFiltered := string(filesCmd)
	cookbookRegexp := regexp.MustCompile(`(.*/cookbooks/[^/]*)/(.*)`)
	filesCmdFiltered := cookbookRegexp.ReplaceAllString(filesCmdUnFiltered, "$1")
	filesList := helpers.RemoveDuplicatesUnordered(strings.Split(strings.TrimSpace(filesCmdFiltered),"\n"))

	return filesList
}
