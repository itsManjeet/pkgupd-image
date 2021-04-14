package utils

import (
	"log"
	"os"
	"os/exec"
)

func Executor(c string, d string, env []string) error {

	environment := os.Environ()
	environment = append(environment, env...)

	cmd := exec.Command("sh", "-ce", c)
	cmd.Dir = d
	cmd.Env = append(cmd.Env, environment...)
	out, err := cmd.CombinedOutput()
	log.Println(string(out))

	if err != nil {
		return err
	}

	return nil
}
