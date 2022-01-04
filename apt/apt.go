package apt

import (
	"compress/gzip"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/itsmanjeet/pkgupd-image/utils"
	"gopkg.in/ini.v1"
)

// Package holds the debian package information
type Package struct {
	Name        string `ini:"Package"`
	Url         string `ini:"Filename"`
	Version     string `ini:"Version"`
	Description string `ini:"Description"`
	Depends     string `ini:"Depends"`
	Provides    string `ini:"Provides"`
}

func (pkg Package) GetProvides() []string {
	provides := make([]string, 0)

	split_provides := strings.Split(pkg.Provides, ",")
	if len(split_provides) == 0 {
		return provides
	}

	for _, prov := range split_provides {
		provide := strings.TrimSpace(prov)
		if len(provide) == 0 {
			continue
		}
		provides = append(provides, provide)
	}

	return provides
}

func (pkg Package) GetDepends() []string {
	depends := []string{}
	for _, dep := range strings.Split(pkg.Depends, ",") {
		start_idx := strings.Index(dep, "(")
		if start_idx != -1 {
			dep = dep[:start_idx]
		}

		start_idx = strings.Index(dep, "|")
		if start_idx != -1 {
			dep = dep[:start_idx]
		}

		start_idx = strings.Index(dep, ":")
		if start_idx != -1 {
			dep = dep[:start_idx]
		}

		depend := strings.TrimSpace(dep)
		if len(depend) == 0 {
			continue
		}
		depends = append(depends, depend)
	}
	return depends
}

type Apt struct {
	WorkDir      string
	Database     []*Package
	Mirror       string
	Repositories []string
	Release      string
	Architecture string
}

func (apt Apt) Get(pkgid string) *Package {

	for _, pkg := range apt.Database {
		if pkg.Name == pkgid {
			return pkg
		}

		for _, prov := range pkg.GetProvides() {
			if prov == pkgid {
				log.Println("using", pkg.Name, "for", pkgid)
				return pkg
			}
		}
	}

	log.Println("Error! " + pkgid + " not in database")
	return nil
}

func (apt *Apt) Sync() error {
	log.Println("syncing database")
	apt.Database = make([]*Package, 0)

	workdir, err := ioutil.TempDir("", "appimage-*")
	if err != nil {
		log.Println("failed to create temporary directory", err)
		return err
	}
	defer os.RemoveAll(workdir)

	database_dir := path.Join(workdir, "database")
	if _, err := os.Stat(database_dir); err != nil {
		if err := os.MkdirAll(database_dir, 0755); err != nil {
			log.Println("Failed to creating database dir", err)
			return err
		}
	}
	for _, repository := range apt.Repositories {
		packageGZ_Url := apt.Mirror + path.Join("dists", apt.Release, repository, "binary-"+apt.Architecture, "Packages.gz")
		log.Println("retreving package from", packageGZ_Url)
		packageGZ_filepath := path.Join(database_dir, repository+".gz")

		if err := utils.Download(packageGZ_filepath, packageGZ_Url); err != nil {
			log.Println("Failed to download package file for", repository, ", Error", err)
			continue
		}

		file, err := os.Open(packageGZ_filepath)
		if err != nil {
			log.Println("failed to read package file for", repository, ", Error", err)
			continue
		}
		defer file.Close()

		gzip_reader, err := gzip.NewReader(file)
		if err != nil {
			log.Println("failed to create gzip reader", err)
			continue
		}
		defer gzip_reader.Close()

		data, err := ioutil.ReadAll(gzip_reader)
		if err != nil {
			log.Println("failed to read gzip conent", err)
			continue
		}

		for _, data_patch := range strings.Split(string(data), "\n\n") {
			ini_file, err := ini.Load([]byte(data_patch))
			if err != nil {
				log.Println("error: ", err)
				continue
			}

			pkg := &Package{}
			if err := ini_file.MapTo(pkg); err != nil {
				log.Println("failed to read map package data", err)
				continue
			}

			apt.Database = append(apt.Database, pkg)
		}
	}

	return nil

}
