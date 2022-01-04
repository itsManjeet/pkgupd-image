package image

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/itsmanjeet/pkgupd-image/recipe"
	"github.com/itsmanjeet/pkgupd-image/utils"
)

type Image struct {
	Directory string
	Recipe    *recipe.Recipe
}

const (
	APPRUN_URL = "https://github.com/AppImage/AppImageKit/releases/download/13/AppRun-x86_64"
)

func (img Image) installAppRun() error {
	appRunPath := path.Join(img.Directory, "AppRun")
	if len(img.Recipe.AppRun) != 0 {
		if err := ioutil.WriteFile(appRunPath, []byte(img.Recipe.AppRun), 0755); err != nil {
			log.Println("failed to install custom AppRun script", err)
			return err
		}
	} else {
		if err := utils.Download(appRunPath, APPRUN_URL); err != nil {
			log.Println("failed to download AppRun from", APPRUN_URL, err)
			return err
		}
	}

	if err := os.Chmod(appRunPath, 0755); err != nil {
		log.Println("failed to chmod AppRun", err)
		return err
	}

	return nil
}

func (img Image) installDesktopFile() error {
	desktopfilePath := path.Join(img.Directory, "app.desktop")
	if len(img.Recipe.Desktop) != 0 {
		if err := ioutil.WriteFile(desktopfilePath, []byte(img.Recipe.Desktop), 0644); err != nil {
			log.Println("failed to install custom desktop file", err)
			return err
		}
	}
	if _, err := os.Stat(desktopfilePath); err != nil {
		log.Println("no desktop file found", err)
		return errors.New("no desktop file found")
	}
	return nil
}
