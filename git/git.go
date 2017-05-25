package git

import (
  //"os"
	"os/exec"
 //"github.com/nlopes/slack"
 "fmt"
 "log"
 "strings"
 "regexp"
 "../helpers"
 "../structs"
)


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

//	fmt.Println("TEST"+ n1["commit"])
		fmt.Println("Commit: "+ commit)
		fmt.Println("Author: "+ author)
    fmt.Println("Date: "+ date)
    fmt.Println("Title: "+ title)
	return structs.CommitInfo{merge, author, date}
}

func FilesListForEachCommit(commit string) []string {
	fmt.Println("Retrieving files list for the commit " + commit )
	filesCmd, err := exec.Command("git", "ls-tree", "-r", "--name-only", commit).Output()
	if err != nil {
		log.Fatal(err)
	}
	filesCmdUnFiltered := string(filesCmd)
	cookbookRegexp := regexp.MustCompile(`(.*/cookbooks/[^/]*)/(.*)`)
	filesCmdFiltered := cookbookRegexp.ReplaceAllString(filesCmdUnFiltered, "$1")
	filesList := helpers.RemoveDuplicatesUnordered(strings.Split(strings.TrimSpace(filesCmdFiltered),"\n"))

	return filesList
}
