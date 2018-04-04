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
func CookbookInfo(path string) (name string, version string, err error) {
	path += "/metadata.rb"
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Cannot read metadata file for cookbook %s", path)
		log.Error(err)
		return name, version, err
	}

	lines := strings.Split(string(input), "\n")

	for _, line := range lines {
		cookbookNamePattern, _ := regexp.MatchString("^( |\t)*name( |\t)+(\\'|\")(.+)(\\'|\")", line)
		cookbookVersionPattern, _ := regexp.MatchString("^( |\t)*version( |\t)+(\\'|\")(.+)(\\'|\")", line)

		if cookbookNamePattern {
			line = strings.Replace(line, "\"", "'", 2)
			lineArray := strings.Split(line, "'")
			name = lineArray[1]
		}

		if cookbookVersionPattern {
			line = strings.Replace(line, "\"", "'", 2)
			lineArray := strings.Split(line, "'")
			version = lineArray[1]
		}
	}

	if name == "" || version == "" {
		log.Errorf("Cookbook name or version not found on the following file %s", path)
		return name, version, err
	}

	return name, version, nil
}
