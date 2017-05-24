package main

import (
	//"github.com/sysbart/spife"
	//"os"
	"os/exec"
 //"github.com/nlopes/slack"
 "fmt"
 "log"
 "strings"
 "regexp"
 "flag"
)


type CommitInfo struct {
	IsMergeCommit bool
	Author     string
	Date     string
}

// https://www.dotnetperls.com/duplicates-go
func removeDuplicatesUnordered(elements []string) []string {
    encountered := map[string]bool{}

    // Create a map of all unique elements.
    for v:= range elements {
        encountered[elements[v]] = true
    }

    // Place all keys from the map into a slice.
    result := []string{}
    for key, _ := range encountered {
        result = append(result, key)
    }
    return result
}

func notifySlack() {

}





func runKnife(cmd string) {
	uploadCmd, err := exec.Command("echo", "knife " + cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.TrimSpace(string(uploadCmd)))
}

func chefUpload(object string, file string) {
	fmt.Println(strings.Title(object) + " " + file + " has been modified")

	if object == "cookbook" {
	//	baseFolder := regexp.MustCompile(`(.*/cookbooks/[^/]*)/(.*)`).
		runKnife(object + " upload -o " + file)
	} else {
		runKnife(object + " from file " + file)
	}

	fmt.Println(strings.Title(object) + " " + file + " uploaded")
}

func getLastCommit() string {
	commitCmd, err := exec.Command("git", "log", "-n1", "--pretty=format:%h").Output()
	if err != nil {
		log.Fatal(err)
	}

	commit := string(commitCmd)
	fmt.Println("Last commit found: " + commit)
	return commit
}

func getCommitInfo(commit string) string {
	commitCmd, err := exec.Command("git", "log", "-n1").Output()
	if err != nil {
		log.Fatal(err)
	}

	commitRegexp := regexp.MustCompile(`commit (?P<commit>.*)|Merge: (.*)|Author: (.*)|Date: (.*)|(.*)`)
	n1 := commitRegexp.SubexpNames()
	r2 := commitRegexp.FindAllStringSubmatch(string(commitCmd), -1)[1][0]

//	fmt.Println("TEST"+ n1["commit"])
		fmt.Println("TEST"+ r2)
	return n1[0]
}

func getFilesListForEachCommit(commit string) []string {
	fmt.Println("Retrieving files list for the commit " + commit )
	filesCmd, err := exec.Command("git", "ls-tree", "-r", "--name-only", commit).Output()
	if err != nil {
		log.Fatal(err)
	}
	filesCmdUnFiltered := string(filesCmd)
	cookbookRegexp := regexp.MustCompile(`(.*/cookbooks/[^/]*)/(.*)`)
	filesCmdFiltered := cookbookRegexp.ReplaceAllString(filesCmdUnFiltered, "$1")
	filesList := removeDuplicatesUnordered(strings.Split(strings.TrimSpace(filesCmdFiltered),"\n"))

	return filesList
}

func main() {
	optCommit := flag.String("commit", "", "commit ID")
	flag.Parse()
	commit := *optCommit
	if commit == "" {
	        commit = getLastCommit()
	}
	files := getFilesListForEachCommit(commit)
	getCommitInfo(commit)
	// Environment and roles are file based, cookbooks are folder based
	for _, file := range files {
		if strings.Contains(file, "environments/") {
			chefUpload("environment", file)
		} else if strings.Contains(file, "roles/") {
			chefUpload("role", file)
		} else if strings.Contains(file, "cookbooks/") {
			chefUpload("cookbook", file)
		}
	}

	//fmt.Println(cookbookFiles)

}
