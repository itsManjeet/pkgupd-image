package utils

import (
	"log"
	"os/exec"
)

func Patch(appdir string) error {
	if len(appdir) == 0 {
		log.Panic("nil value of appdir:", appdir)
	}
	cmd := exec.Command("find", "usr/", "-type", "f", "-exec", "sed", "-i", "-e", "s|/usr|././|g", "{}", "+")
	cmd.Dir = appdir
	out, err := cmd.CombinedOutput()
	log.Println(string(out))
	if err != nil {
		return err
	}
	cmd = exec.Command("find", "usr/", "-type", "f", "-exec", "sed", "-i", "-e", "s@././/bin/env@/usr/bin/env@g", "{}", "+")
	cmd.Dir = appdir
	out, err = cmd.CombinedOutput()
	log.Println(string(out))

	// cmd = exec.Command("bash", "-c", "cp -ra usr/* .")
	// cmd.Dir = appdir
	// log.Println(cmd.Run())

	// os.RemoveAll(appdir + "/usr/")

	return err
}
