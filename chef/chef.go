package chef

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func ParseObjectByFileName(object string, file string) string {
	var parsedFilename string
	if object == "cookbook" {
		parsedFilename = CookbookInfo(file)
	} else {
		fileRegexp := regexp.MustCompile(`.*/(.*)\..*$`)
		parsedFilename = fileRegexp.ReplaceAllString(file, "$1")
	}
	return parsedFilename
}

func Upload(object string, file string) {
	fmt.Printf("%s %s has been modified\n", strings.Title(object), file)
	parsedFilename := ParseObjectByFileName(object, file)

	if object == "cookbook" {
		cookbookRegexp := regexp.MustCompile(`(.*cookbooks)/(.*)`)
		cookbookBaseFolder := cookbookRegexp.ReplaceAllString(file, "$1")
		knife(object, "upload", "-o", cookbookBaseFolder, parsedFilename)
	} else {
		knife(object, "from", "file", file)
	}

	fmt.Printf("%s %s has been uploaded to the Chef server\n", strings.Title(object), parsedFilename)
}

func knife(cmd ...string) {
	uploadCmd, err := exec.Command("knife", cmd...).Output()
	if err != nil {
		log.Print(string(uploadCmd))
		log.Fatal(err)
	}
	fmt.Println(strings.TrimSpace(string(uploadCmd)))
}

func CookbookInfo(path string) string {
	path += "/metadata.rb"
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	var cookbookName string

	lines := strings.Split(string(input), "\n")

	for _, line := range lines {
		cookbookNamePattern, _ := regexp.MatchString("^( |\t)*name( |\t)+(\\'|\")(.+)(\\'|\")", line)

		if cookbookNamePattern {
			line = strings.Replace(line, "\"", "'", 2)
			lineArray := strings.Split(line, "'")
			cookbookName = lineArray[1]
		}
	}

	if cookbookName == "" {
		log.Fatalln("Cookbook name not found on the following file : " + path)
	}

	return cookbookName
}
