package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App struct {
		ID      string `yaml:"id"`
		Version string `yaml:"version"`
		About   string `yaml:"desc"`
	} `yaml:"app"`

	Union  bool `yaml:"union"`
	Patch  bool `yaml:"patch"`
	Distro struct {
		ID           string   `yaml:"id"`
		Version      string   `yaml:"version"`
		Mirror       string   `yaml:"mirror"`
		Repositories []string `yaml:"repo"`
		Skips        []string `yaml:"skips"`
		Includes     []string `yaml:"includes"`
	} `yaml:"distro"`

	Execute struct {
		Sources     []string `yaml:"sources"`
		Script      string   `yaml:"script"`
		Environment []string `yaml:"environ"`
	} `yaml:"exec"`

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
