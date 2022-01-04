package apt

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/itsManjeet/app-fctry/utils"
)

func (apt Apt) Install(rootdir string, pkg_ids ...string) error {
	log.Println("installing", pkg_ids)
	pkgs := apt.packagesList(pkg_ids...)
	if pkgs == nil {
		return errors.New("failed to resolve depends")
	}

	pkgs_dir, err := ioutil.TempDir("", "pkgupd-apps-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(pkgs_dir)

	extractFile := func(tarfile string) error {
		return exec.Command("bsdtar", "-xf", tarfile, "-C", rootdir).Run()
	}

	for _, pkg := range pkgs {

		if len(pkg.Url) == 0 {
			log.Println("no url for " + pkg.Name)
			continue
		}

		pkg_url := apt.Mirror + pkg.Url
		pkg_path := path.Join(pkgs_dir, path.Base(pkg.Url))

		if err := utils.Download(pkg_path, pkg_url); err != nil {
			return err
		}

		log.Println("extracting", pkg_path)
		if err := extractFile(pkg_path); err != nil {
			return err
		}

		log.Println("extracting data", pkg_path)

		data_found := false
		for _, ext := range []string{".xz", ".gz"} {
			datafile := path.Join(rootdir, "data.tar"+ext)
			if _, err := os.Stat(datafile); err == nil {
				data_found = true
				if err := extractFile(datafile); err != nil {
					return err
				}
			}
		}
		if !data_found {
			log.Println("no data file found")
			return errors.New("no data file found in" + pkg.Name)
		}

	}

	for _, cache := range []string{"data.tar.xz", "control.tar.xz", "debian-binary", "control.tar.gz", "data.tar.gz"} {
		os.RemoveAll(path.Join(rootdir, cache))
	}

	return nil
}
