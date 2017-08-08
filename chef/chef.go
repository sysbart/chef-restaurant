package chef

import (
	"io/ioutil"
	"regexp"
	"strings"
	log "github.com/sirupsen/logrus"
	"github.com/sysbart/chef-restaurant/helpers"
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
	log.Infof("%s %s has been modified", strings.Title(object), file)
	parsedFilename := ParseObjectByFileName(object, file)

	if object == "cookbook" {
		cookbookRegexp := regexp.MustCompile(`(.*cookbooks)/(.*)`)
		cookbookBaseFolder := cookbookRegexp.ReplaceAllString(file, "$1")
		knife(object, "upload", "-o", cookbookBaseFolder, parsedFilename)
	} else if strings.HasSuffix(file, ".json") || strings.HasSuffix(file, ".rb") {
		knife(object, "from", "file", file)
	} else {
		log.Infof("The file %s has been not uploaded since its filetype is not supported", file)
	}

	log.Infof("%s %s has been uploaded to the Chef server\n", strings.Title(object), parsedFilename)
}

func knife(cmd ...string) {
	helpers.RunCommand("knife", cmd...)
}

func CookbookInfo(path string) string {
	path += "/metadata.rb"
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal("Cookbook name not found on the following file : " + path)
	}

	return cookbookName
}
