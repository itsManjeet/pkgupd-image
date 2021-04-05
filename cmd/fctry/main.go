package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s /path/to/recipe.yml\n", os.Args[0])
		os.Exit(1)
	}

	basedir := "/tmp/"
	if len(os.Args) > 2 {
		basedir = os.Args[2]
	}

	factory := Factory{
		recipefile: os.Args[1],
		basedir:    basedir,
	}

	if err := factory.Build(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(":: Build Success ::")
}
