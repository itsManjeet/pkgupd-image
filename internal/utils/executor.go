package utils

import (
	"log"
	"os"
	"os/exec"
)

func Executor(cmds []string, d string, env []string) error {

	environment := os.Environ()
	environment = append(environment, env...)

	for _, c := range cmds {
		cmd := exec.Command("sh", "-c", c)
		cmd.Dir = d
		cmd.Env = append(cmd.Env, environment...)
		out, err := cmd.CombinedOutput()
		log.Println(string(out))

		if err != nil {
			return err
		}

		environment = append(environment, cmd.Env...)
	}

	return nil
}
