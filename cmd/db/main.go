package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/itsManjeet/app-fctry/config"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Printf("\nUsage: %s path/to/recipe out.db", args[0])
		os.Exit(1)
	}

	recipeDir := args[1]
	output := args[2]
	dirData, err := ioutil.ReadDir(recipeDir)
	if err != nil {
		log.Println("Error!", err)
		os.Exit(1)
	}

	recipeData := make([]config.Config, 0)
	for _, file := range dirData {
		recipe, err := config.Load(path.Join(recipeDir, file.Name()))
		if err != nil {
			log.Panicln("Panic!", file.Name(), err)
		}

		recipeData = append(recipeData, *recipe)
	}

	data, err := json.Marshal(recipeData)
	if err != nil {
		log.Panicln("Panic!", err)
	}

	if err := ioutil.WriteFile(output, data, 0755); err != nil {
		log.Panicln("Panic!", err)
	}
}
