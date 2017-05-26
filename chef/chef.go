package chef

import (
 "os/exec"
 "fmt"
 "log"
 "strings"
 "io/ioutil"
 "regexp"
)

func Upload(object string, file string) {
	fmt.Println(strings.Title(object) + " " + file + " has been modified")

	if object == "cookbook" {
	//	baseFolder := regexp.MustCompile(`(.*/cookbooks/[^/]*)/(.*)`).
		knife(object + " upload -o " + file)
	} else {
		knife(object + " from file " + file)
	}

	fmt.Println(strings.Title(object) + " " + file + " uploaded")
}

func knife(cmd string) {
	uploadCmd, err := exec.Command("echo", "knife " + cmd).Output()
	if err != nil {
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
    log.Fatalln("Cookbook name not found on the following file : " + path )
  }

  return cookbookName
}

/*
func CookbookGenerateNewVersion(folder string) []string {

}

func CookbookVersionType() string {

}

func CookbookModifyMetadata(folder string, newVersion string) {

} */
