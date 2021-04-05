package utils

import (
	"log"
	"os"
	"os/exec"
)

func Patch(appdir string) error {
	cmd := exec.Command("find", "usr/", "-type", "f", "-executable", "-exec", "sed", "-i", "-e", "s|/usr|././|g", "{}", "+")
	cmd.Dir = appdir
	out, err := cmd.CombinedOutput()
	log.Println(string(out))

	cmd = exec.Command("bash", "-c", "cp -ra usr/* .")
	cmd.Dir = appdir
	log.Println(cmd.Run())

	os.RemoveAll(appdir + "/usr/")
	return err
}
