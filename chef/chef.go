package chef

import (
 "os/exec"
 "fmt"
 "log"
 "strings"
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
