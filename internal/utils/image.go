package utils

import (
	"log"
	"os/exec"
)

func MakeImage(appdir, outdir string) error {
	cmd := "ARCH=x86_64 appimagetool -v -n " + appdir + " " + outdir
	log.Println("executing: ", cmd)
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	log.Println(string(out))
	return err
}
