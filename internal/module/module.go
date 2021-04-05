package module

import (
	"appfctry/internal/config"
	"appfctry/internal/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type Package struct {
	Name, URL   string
	Version     string
	Description string
	Depends     []string
}

type Module interface {
	GetURL(string, string, string) string
	Prepare(string) (map[string]Package, error)
	Install(string, string) error
}

type Executor struct {
	module      Module
	config      *config.Config
	database    map[string]Package
	applist     []Package
	alreadydone []string
}

func Initialize(c *config.Config, m Module) *Executor {
	return &Executor{
		config: c,
		module: m,
	}
}

func (exec Executor) GetApp(appID string) (Package, error) {
	if app, ok := exec.database[appID]; ok {
		return app, nil
	}

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
	exec.alreadydone = exec.config.Skip
	exec.applist = make([]Package, 0)

	if _, err := os.Stat("assets/apps.list"); err == nil {
		if data, err := ioutil.ReadFile("assets/apps.list"); err == nil {
			appslist := strings.Split(string(data), "\n")
			log.Println("added", len(appslist), "to skip list")
			exec.config.Skip = append(exec.config.Skip, appslist...)
		}
	}

	for _, repo := range exec.config.Repositories {
		log.Println("syncing ", repo)

		url := exec.module.GetURL(exec.config.URL, exec.config.Version, repo)
		filepath := wdir + "/" + repo

		if err := utils.DownloadFile(filepath, url); err != nil {
			return err
		}
		db, err := exec.module.Prepare(filepath)
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

	for _, d := range exec.config.Include {
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

	if err := utils.DownloadFile(filepath, exec.config.URL+"/"+app.URL); err != nil {
		return err
	}

	if err := exec.module.Install(filepath, wrkdir); err != nil {
		return err
	}

	return nil
}
