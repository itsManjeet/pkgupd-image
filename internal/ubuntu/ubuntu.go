package ubuntu

import (
	"appfctry/internal/module"
	"appfctry/internal/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Ubuntu struct {
}

func (d *Ubuntu) GetURL(mirr, repo, version string) string {
	return fmt.Sprintf("%s/dists/%s/%s/binary-amd64/Packages.gz", mirr, repo, version)
}

func (d *Ubuntu) Prepare(path string) (map[string]module.Package, error) {

	output, err := d.readfile(path)
	if err != nil {
		return nil, err
	}

	db := make(map[string]module.Package)

	lines := strings.Split(output, "\n")
	i := 0

	getval := func(l string) (string, string) {

		if strings.Contains(l, ":") {
			idx := strings.Index(l, ":")
			return strings.TrimSpace(l[:idx]), strings.TrimSpace(l[idx+1:])
		}
		return "", ""
	}

	getdep := func(v string) []string {
		deps := make([]string, 0)
		dep := strings.Split(v, ",")
		for _, d := range dep {
			dv := strings.Split(strings.TrimSpace(d), " ")
			dv = strings.Split(dv[0], ":")
			deps = append(deps, dv[0])
		}
		return deps
	}

	for i <= len(lines) {
		var pkg module.Package
		pkg.Depends = make([]string, 0)

		for {
			t, v := getval(lines[i])
			switch t {
			case "Package":
				pkg.Name = v
			case "Version":
				pkg.Version = v
			case "Description":
				pkg.Description = v
			case "Depends":
				pkg.Depends = append(pkg.Depends, getdep(v)...)

			case "Filename":
				pkg.URL = v
			}
			i++

			if i >= len(lines) || strings.HasPrefix(lines[i], "Package") {
				db[pkg.Name] = pkg
				break
			}
		}

		if i >= len(lines) {
			break
		}
	}

	return db, nil
}

func (d *Ubuntu) readfile(path string) (string, error) {
	data, err := exec.Command("gzip", "--stdout", "--decompress", path).Output()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (d *Ubuntu) Install(file, dir string) error {
	fmt.Println(":: Extracting", file, "::")
	if err := utils.Extractfile(file, dir); err != nil {
		return err
	}

	if err := utils.Extractfile(dir+"/data.tar.gz", dir); err != nil {
		if err := utils.Extractfile(dir+"/data.tar.xz", dir); err != nil {
			return err
		}
	}

	for _, cache := range []string{"data.tar.xz", "control.tar.xz", "debian-binary", "control.tar.gz", "data.tar.gz"} {
		os.RemoveAll(dir + "/" + cache)
	}
	return nil
}
