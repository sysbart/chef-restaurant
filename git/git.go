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

func CommitInfo(commit string) *structs.CommitInfo {
	commitCmd, err := exec.Command("git", "log", "-n1").Output()
	if err != nil {
		log.Fatal(err)
	}

	commitRegexp := regexp.MustCompile(`commit (?P<commit>.*)|Merge: (.*)|Author: (.*)|Date: (.*)|(.*)`)
	n1 := commitRegexp.SubexpNames()
	r2 := commitRegexp.FindAllStringSubmatch(string(commitCmd), -1)[1][0]

//	fmt.Println("TEST"+ n1["commit"])
		fmt.Println("TEST"+ r2)
	return false, "b", "c"
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
