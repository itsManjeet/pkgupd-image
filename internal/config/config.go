package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App string `yaml:"app"`

	Module  string `yaml:"module"`
	Version string `yaml:"version"`

	Patch bool `yaml:"patch"`

	URL          string   `yaml:"url"`
	Repositories []string `yaml:"repositories"`
	Script       []string `yaml:"script"`
	Environment  []string `yaml:"environment"`
	Include      []string `yaml:"include"`
	Skip         []string `yaml:"skip"`

	AppRun  string `yaml:"appRun"`
	Desktop string `yaml:"desktop"`
}

func Load(filename string) (*Config, error) {
	var config Config

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
