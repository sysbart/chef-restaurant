package helpers
import (
	"os/exec"
  log "github.com/sirupsen/logrus"
  "strings"
)

var WorkingFolder string

// Based on https://www.dotnetperls.com/duplicates-go
func RemoveDuplicatesUnordered(elements []string) []string {
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


func RunCommand(cmd string, args...string) []byte {
  runCmd := exec.Command(cmd, args...)
	runCmd.Dir = WorkingFolder
	outputCmd, err := runCmd.Output()
  log.Debugf("Running command : %s %s", cmd, strings.Join(args[:], " "))
  log.Debugf("Using the following working directory : %s", WorkingFolder)

  if err != nil {
    log.Fatal(err)
  }
  return outputCmd
}
