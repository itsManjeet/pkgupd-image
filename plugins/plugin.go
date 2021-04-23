package plugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/itsManjeet/app-fctry/config"
	"github.com/itsManjeet/app-fctry/utils"
)

type Package struct {
	Name, URL   string
	Version     string
	Description string
	Depends     []string
	Provides    []string
}

type Plugin interface {
	GetURL(string, string, string) string
	Prepare(string) (map[string]Package, error)
	Install(string, string) error
}

type Executor struct {
	plugin      Plugin
	config      *config.Config
	database    map[string]Package
	applist     []Package
	alreadydone []string
}

func Initialize(c *config.Config, m Plugin) *Executor {
	return &Executor{
		config: c,
		plugin: m,
	}
}

func (exec Executor) GetApp(appID string) (Package, error) {
	if app, ok := exec.database[appID]; ok {
		return app, nil
	}

	for _, app := range exec.database {
		for _, a := range app.Provides {
			if appID == a {
				return app, nil
			}
		}
	}

	fmt.Println("Error! " + appID + " not found in repo")
	os.Exit(1)
	return Package{}, nil
}

func (exec *Executor) inlist(appID Package) bool {
	for _, a := range exec.applist {
		if a.Name == appID.Name {
			return true
		}
	}
	return false
}

func (exec *Executor) isdone(appID string) bool {
	for _, a := range exec.alreadydone {
		if a == appID {
			return true
		}
	}
	return false
}

func (exec *Executor) caldep(app Package) {
	for _, a := range app.Depends {
		if exec.isdone(a) {
			continue
		}

		exec.alreadydone = append(exec.alreadydone, a)
		apd, err := exec.GetApp(a)
		if err != nil {
			fmt.Println(apd, "is missing required by", app.Name)
			os.Exit(1)
		}
		exec.caldep(apd)
		if !exec.inlist(apd) {
			exec.applist = append(exec.applist, apd)
		}
	}

	if !exec.inlist(app) {
		exec.applist = append(exec.applist, app)
	}
}

func (exec *Executor) mergedb(db map[string]Package) {
	for p, a := range db {
		exec.database[p] = a
	}
}

func (exec *Executor) Sync(wdir string) error {
	fmt.Println(":: Syncing database ::")
	exec.database = make(map[string]Package)
	exec.alreadydone = exec.config.Distro.Skips
	exec.applist = make([]Package, 0)

	if _, err := os.Stat("assets/apps.list"); err == nil {
		if data, err := ioutil.ReadFile("assets/apps.list"); err == nil {
			appslist := strings.Split(string(data), "\n")
			log.Println("added", len(appslist), "to skip list")
			exec.config.Distro.Skips = append(exec.config.Distro.Skips, appslist...)
		}
	}

	for _, repo := range exec.config.Distro.Repositories {
		log.Println("syncing ", repo)

		url := exec.plugin.GetURL(exec.config.Distro.Mirror, exec.config.Distro.Version, repo)
		filepath := wdir + "/" + repo

		if err := utils.DownloadFile(filepath, url); err != nil {
			return err
		}
		db, err := exec.plugin.Prepare(filepath)
		if err != nil {
			return err
		}

		exec.mergedb(db)
	}

	return nil
}

func (exec *Executor) Depends(appID string) []Package {
	app, err := exec.GetApp(appID)
	if err != nil {
		fmt.Println(appID, "not found in db")
		os.Exit(1)
	}

	exec.caldep(app)

	for _, d := range exec.config.Distro.Includes {
		apd, err := exec.GetApp(d)
		if err != nil {
			fmt.Println(d, "not found in db")
			os.Exit(1)
		}
		exec.caldep(apd)
	}

	return exec.applist
}

func (exec *Executor) Install(appID string, srcdir, wrkdir string) error {
	if len(appID) == 0 {
		return nil
	}
	fmt.Println(":: Installing", appID, "to", wrkdir, "::")

	app, _ := exec.GetApp(appID)
	if len(app.Name) == 0 {
		fmt.Println("Error! no package with name", appID)
		return errors.New("No package with id " + appID)
	}

	_, file := path.Split(app.URL)
	filepath := srcdir + "/" + file

	if err := utils.DownloadFile(filepath, exec.config.Distro.Mirror+"/"+app.URL); err != nil {
		return err
	}

	if err := exec.plugin.Install(filepath, wrkdir); err != nil {
		return err
	}

	return nil
}
