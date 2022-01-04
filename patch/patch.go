package patch

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Patch struct {
	Directory string
}

// UsrRelocatable patch /usr to ././ in ./usr
// https://github.com/AppImage/pkg2appimage/blob/master/functions.sh#L77
// Make content usr/ relocatable
func (patch Patch) UsrRelocateable() error {

	script := `
	find usr/ -type f -exec sed -i -e 's|/usr|././|g' {} \;
  find usr/ -type f -exec sed -i -e 's@././/bin/env@/usr/bin/env@g' {} \;
	`

	exec_command := func(bin string, args ...string) error {
		cmd := exec.Command(bin, args...)
		cmd.Dir = patch.Directory

		return cmd.Run()
	}
	if err := exec_command("bash", "-c", script); err != nil {
		log.Println("failed to patch usr", err)
		return err
	}

	if err := exec_command("rsync", "-a", "-v", "usr/", "."); err != nil {
		log.Println("failed to merge usr", err)
		return err
	}

	if err := os.RemoveAll(path.Join(patch.Directory, "usr")); err != nil {
		log.Println("failed to clean usr", err)
		return err
	}

	return nil
}

// Cleanup clean unneccessary cache
// like /usr/include, /lib/cmake, /lib/pkgconfig
func (patch Patch) Cleanup() error {
	for _, lib := range strings.Split(EXCLUE_LIBS, "\n") {
		if len(lib) == 0 {
			continue
		}
		for _, libdir := range []string{"lib", "lib/x86_64-linux-gnu", "usr/lib", "usr/lib/x86_64-linux-gnu"} {
			libpath := path.Join(patch.Directory, libdir, lib)
			if _, err := os.Stat(libpath); err == nil {
				log.Println("found and cleaning", libpath)
				if err := os.Remove(libpath); err != nil {
					log.Println("failed to remove", err)
					break
				}
			}
		}
	}
	for _, dir := range []string{
		"usr/include", "usr/lib/cmake", "usr/lib/pkgconfig",
	} {
		dir = path.Join(patch.Directory, dir)
		if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
			if err := os.RemoveAll(dir); err != nil {
				log.Println("failed to clean", dir)
				return err
			}
		}
	}
	return nil
}
