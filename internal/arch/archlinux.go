package archlinux

import (
	"appfctry/internal/config"
	"appfctry/internal/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type ArchLinux struct {
	Config    config.Config
	Database  map[string]map[string]ArchPkg
	Installed []string
}

func (a ArchLinux) geturl(pkgname string) (string, error) {
	pkg, repo, err := a.getpkg(pkgname)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/os/x86_64/%s", a.Config.ArchLinux.Mirror, repo, pkg.Filename), nil
}

func (a ArchLinux) getpkg(pkgname string) (ArchPkg, string, error) {
	if strings.Contains(pkgname, ">=") {
		pkgname = strings.Split(pkgname, ">=")[0]
	} else if strings.Contains(pkgname, "=") {
		pkgname = strings.Split(pkgname, "=")[0]
	}

	for repo, db := range a.Database {
		if val, ok := db[pkgname]; ok {
			return val, repo, nil
		}
	}

	for repo, db := range a.Database {
		for _, data := range db {
			for _, a := range data.Provides {
				if pkgname == a {
					return data, repo, nil
				}
			}
		}
	}

	return ArchPkg{}, "", errors.New("No package found for name " + pkgname)
}

func (a *ArchLinux) downloadDatabase() error {
	a.Database = make(map[string]map[string]ArchPkg)
	for _, i := range a.Config.ArchLinux.Repositories {
		fmt.Println("syncing", i)
		dataurl := fmt.Sprintf("%s/%s/os/x86_64/%s.db.tar.gz", a.Config.ArchLinux.Mirror, i, i)
		outfile := a.Config.SRCDir + "/" + i + ".tar.gz"
		if _, err := os.Stat(outfile); os.IsNotExist(err) {
			log.Println("using url", dataurl)
			if err := utils.DownloadFile(outfile, dataurl); err != nil {
				log.Println("failed to sync", i)
				return err
			}
		}

		output, err := exec.Command("tar", "-xaf", outfile, "-O").Output()
		if err != nil {
			return err
		}

		a.Database[i], err = a.parseDatabase(string(output))
		if err != nil {
			return err
		}

	}

	return nil
}

func (a ArchLinux) Install(pkgname string) error {

	pkg, _, err := a.getpkg(pkgname)
	if err != nil {
		return err
	}

	pkgname = pkg.Name

	for _, a := range a.Installed {
		if a == pkgname {
			return nil
		}
	}

	srcdir := a.Config.SRCDir
	wrkdir := a.Config.WRKDir

	os.MkdirAll(wrkdir, 0755)

	for _, i := range pkg.Depends {
		if err := a.Install(i); err != nil {
			return err
		}
	}

	url, _ := a.geturl(pkgname)

	log.Println("downloading file ", url)
	err = utils.DownloadFile(srcdir+"/"+pkg.Filename, url)
	if err != nil {
		log.Println("failed to download bash")
		return err
	}

	log.Println("Extracting file", pkg.Filename)
	if err := utils.Extractfile(srcdir+"/"+pkg.Filename, wrkdir); err != nil {
		log.Println("failed to extract bash")
		return err
	}

	for _, i := range []string{
		".MTREE",
		".BUILDINFO",
		".PKGINFO",
	} {
		os.Remove(wrkdir + "/" + i)
	}

	a.Installed = append(a.Installed, pkgname)

	return nil
}

func (a ArchLinux) Build() error {
	log.Println("syncing database")
	a.Installed = a.Config.ArchLinux.SkipPackage
	a.Installed = append(a.Installed, a.Config.SKIPPACKAGE...)
	if err := a.downloadDatabase(); err != nil {
		return err
	}

	for _, dep := range a.Config.Include {
		if err := a.Install(dep); err != nil {
			return err
		}
	}

	return a.Install(a.Config.App)
}
