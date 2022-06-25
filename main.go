package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/itsmanjeet/pkgupd-image/apprun"
	"github.com/itsmanjeet/pkgupd-image/apt"
	"github.com/itsmanjeet/pkgupd-image/patch"
	"github.com/itsmanjeet/pkgupd-image/union"
)

type arrayFlag []string

func (a *arrayFlag) String() string {
	return "string array representation"
}

func (a *arrayFlag) Set(value string) error {
	*a = append(*a, value)
	return nil
}

var (
	mirror       *string
	release      *string
	arch         *string
	repositories arrayFlag
)

func init() {
	mirror = flag.String("mirror", "http://archive.ubuntu.com/ubuntu/", "Specify custom ubuntu mirror")
	release = flag.String("release", "focal", "Specify Ubuntu release")
	arch = flag.String("arch", "amd64", "Specify Architecure")

	flag.Var(&repositories, "repositories", "Specify repositories")
	flag.Parse()
}

func main() {

	if len(os.Args) > 2 {
		fmt.Printf("Usage: %s [generate|patch|cleanup|union] <args>", os.Args[0])
		os.Exit(1)
	}
	if len(repositories) == 0 {
		repositories = []string{"main", "universe"}
	}

	WorkDir := os.Args[2]
	os.MkdirAll(WorkDir, 0755)

	patcher := patch.Patch{
		Directory: WorkDir,
	}

	task := os.Args[1]
	switch task {
	case "generate":
		apt := apt.Apt{
			Mirror:       *mirror,
			Repositories: repositories,
			Architecture: *arch,
			Release:      *release,
			WorkDir:      WorkDir,
		}

		if err := apt.Sync(); err != nil {
			log.Println("error failed to sync", err)
			os.Exit(1)
		}

		if err := apt.Install(WorkDir, os.Args[3:]...); err != nil {
			log.Println("error failed to build AppDir", err)
			os.Exit(1)
		}

	case "patch":

		log.Println("applying usrpatch")
		if err := patcher.UsrRelocateable(); err != nil {
			log.Println("failed UsrRelocatable() patch", err)
			os.Exit(1)
		}

	case "cleanup":
		if err := patcher.Cleanup(); err != nil {
			log.Println("failed Cleanup() patch", err)
			os.Exit(1)
		}

	case "union":
		if err := union.Install(WorkDir); err != nil {
			log.Println("failed to install libunionpreload", err)
			os.Exit(1)
		}

	case "apprun":
		if err := apprun.Install(WorkDir); err != nil {
			log.Println("failed to install libunionpreload", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Usage: %s [generate|patch] <args>", os.Args[0])
		os.Exit(1)
	}

}
