package main

import (
	"appfctry/internal/archlinux"
	"appfctry/internal/config"
	"appfctry/internal/debian"
	"appfctry/internal/module"
	"appfctry/internal/ubuntu"
	"appfctry/internal/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type Factory struct {
	recipefile string
	basedir    string
	config     *config.Config

	srcdir, pkgdir, wrkdir, syncdir string
}

func contains(list []string, data string) bool {
	for _, a := range list {
		if a == data {
			return true
		}
	}
	return false
}

func (f *Factory) Build() (err error) {

	f.srcdir = f.basedir + "/src/"
	f.pkgdir = f.basedir + "/pkg/"
	f.syncdir = f.basedir + "/sync/"

	f.config, err = config.Load(f.recipefile)
	if err != nil {
		return err
	}

	f.syncdir += f.config.Distro.ID + "/"

	f.wrkdir = f.basedir + "/wrk/" + f.config.App.ID
	f.clean()

	mod, err := f.getModule(f.config.Distro.ID)
	if err != nil {
		return err
	}

	if err := f.pre(); err != nil {
		return err
	}

	appID := f.config.App.ID

	if mod == nil {
		log.Println("Using basic script")
		for _, source := range f.config.Execute.Sources {
			_, file := path.Split(source)
			if err := utils.DownloadFile(f.srcdir+"/"+file, source); err != nil {
				return err
			}

			if err := utils.Extractfile(f.srcdir+"/"+file, f.wrkdir); err != nil {
				return err
			}
		}

	} else {

		exec := module.Initialize(f.config, mod)

		if _, err := os.Stat("assets/apps.list"); err == nil {
			if data, err := ioutil.ReadFile("assets/apps.list"); err == nil {
				appslist := strings.Split(string(data), "\n")
				f.config.Distro.Skips = append(f.config.Distro.Skips, appslist...)
			}
		}

		if err := exec.Sync(f.syncdir); err != nil {
			return err
		}

		if _, err := exec.GetApp(appID); err != nil {
			return err
		}

		for _, dep := range exec.Depends(appID) {
			if contains(f.config.Distro.Skips, dep.Name) {
				continue
			}
			if err := exec.Install(dep.Name, f.srcdir, f.wrkdir); err != nil {
				return err
			}
		}
	}

	if err := utils.Executor(f.config.Execute.Script, f.wrkdir, f.config.Execute.Environment); err != nil {
		return err
	}

	if f.config.Patch {
		fmt.Println("=> Patching image")
		if err := utils.Patch(f.wrkdir); err != nil {
			return err
		}
	}

	icofile := "assets/package.png"
	if _, err := os.Stat(f.wrkdir + "/" + appID + ".png"); os.IsNotExist(err) {
		utils.Copyfile(icofile, f.wrkdir+"/"+appID+".png")
	}

	if err := utils.Copyfile("assets/libunionpreload.so", f.wrkdir+"/libunionpreload.so"); err != nil {
		return err
	}

	utils.Copyfile(f.wrkdir+"/"+appID+".png", f.pkgdir+"/.icons/"+appID+".png")
	desktopfile := f.wrkdir + "/share/applications/" + appID + ".desktop"
	if _, err := os.Stat(desktopfile); err == nil {
		utils.Copyfile(desktopfile, f.wrkdir+"/"+appID+".desktop")
	} else {
		if err := utils.WriteDesktop(f.config.Desktop, appID, f.wrkdir); err != nil {
			return err
		}
	}

	if err := utils.WriteAppRun(f.config.AppRun, appID, f.wrkdir); err != nil {
		return err
	}

	if _, err := os.Stat("assets/files.list"); err == nil {
		if data, err := ioutil.ReadFile("assets/files.list"); err == nil {
			fileslist := strings.Split(string(data), "\n")
			utils.Clean(f.wrkdir, fileslist)
		}
	}

	if err := utils.MakeImage(f.wrkdir, f.pkgdir+"/"+appID); err != nil {
		return err
	}

	f.clean()

	return nil
}

func (f Factory) getModule(module string) (module.Module, error) {
	switch module {
	case "debian":
		return &debian.Debian{}, nil
	case "ubuntu":
		return &ubuntu.Ubuntu{}, nil
	case "archlinux":
		return &archlinux.ArchLinux{}, nil
	}
	if len(module) == 0 {
		return nil, nil
	}

	return nil, errors.New("unsupported module " + module)
}

func (f Factory) pre() error {
	fmt.Println(":: Setting Work Environment ::")

	for _, dir := range []string{
		f.srcdir, f.wrkdir, f.pkgdir, f.syncdir,
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("failed to set %s directory\n", dir)
			return err
		}
	}

	return nil
}

func (f Factory) clean() {
	log.Println("clearing", f.wrkdir)
	if os.Getenv("NO_CLEAN") != "1" {
		os.RemoveAll(f.wrkdir)
	}
}
