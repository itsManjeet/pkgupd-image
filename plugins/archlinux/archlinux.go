package archlinux

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	plugin "github.com/itsManjeet/app-fctry/plugins"
	"github.com/itsManjeet/app-fctry/utils"
)

type ArchLinux struct {
}

func (a *ArchLinux) GetURL(mirr, _, repo string) string {
	return fmt.Sprintf("%s/%s/os/x86_64/%s.db.tar.gz", mirr, repo, repo)
}

func (a ArchLinux) readfile(path string) (string, error) {
	data, err := exec.Command("tar", "-xaf", path, "-O").Output()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (a *ArchLinux) Prepare(p string) (map[string]plugin.Package, error) {

	_, repo := path.Split(p)

	out, err := a.readfile(p)
	if err != nil {
		return nil, err
	}

	db := make(map[string]plugin.Package)

	striper := func(f string) string {
		if strings.Contains(f, ">=") {
			return strings.Split(f, ">=")[0]
		}
		if strings.Contains(f, "=") {
			return strings.Split(f, "=")[0]
		}
		return f
	}

	lines := strings.Split(out, "\n")
	i := 0

	for lines[i] == "%FILENAME%" {
		i++

		var t plugin.Package
		t.URL = repo + "/os/x86_64/" + lines[i]
		i++

		for lines[i] != "%FILENAME%" {
			switch lines[i] {
			case "%NAME%":
				i++
				t.Name = lines[i]

			case "%VERSION%":
				i++
				t.Version = lines[i]

			case "%DEPENDS%":
				i++
				t.Depends = make([]string, 0)
				for lines[i] != "" && lines[i][0] != '%' {
					t.Depends = append(t.Depends, striper(lines[i]))
					i++
				}
			case "%PROVIDES%":
				i++
				t.Provides = make([]string, 0)
				for lines[i] != "" && lines[i][0] != '%' {
					t.Provides = append(t.Provides, striper(lines[i]))
					i++
				}
			}
			i++
			if i >= len(lines) {
				break
			}

		}
		db[t.Name] = t
		if i >= len(lines) {
			break
		}

	}

	return db, nil
}

func (a *ArchLinux) Install(file, dir string) error {
	fmt.Println(":: Extracting", file, "::")
	if err := utils.Extractfile(file, dir); err != nil {
		return err
	}

	for _, cache := range []string{".PKGINFO", ".MTREE"} {
		os.RemoveAll(dir + "/" + cache)
	}
	return nil
}
