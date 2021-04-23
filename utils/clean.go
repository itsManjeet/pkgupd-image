package utils

import (
	"log"
	"os/exec"
)

func Clean(appdir string, files []string) error {
	for _, f := range files {
		if len(f) == 0 {
			continue
		}
		cmd := exec.Command("find", ".", "-name", "*"+f+"*", "-delete", "-print")
		cmd.Dir = appdir
		out, _ := cmd.CombinedOutput()
		log.Println(string(out))
	}

	return nil
}
