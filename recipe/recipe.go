package recipe

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Recipe holds the build information of package
type Recipe struct {
	ID      string `yaml:"id"`
	Version string `yaml:"version"`
	About   string `yaml:"about"`

	Union bool `yaml:"union"`
	Patch bool `yaml:"patch"`

	Mirror       string   `yaml:"mirror"`
	Repositories []string `yaml:"repositories"`
	Includes     []string `yaml:"includes"`
	Release      string   `yaml:"release"`
	Architecture string   `yaml:"arch"`

	Script  string `yaml:"script"`
	AppRun  string `yaml:"AppRun"`
	Desktop string `yaml:"Desktop"`
}

func Load(loc string) (*Recipe, error) {

	data, err := ioutil.ReadFile(loc)
	if err != nil {
		return nil, err
	}

	rcp := &Recipe{}
	err = yaml.Unmarshal(data, rcp)
	if err != nil {
		return nil, err
	}

	return rcp, nil
}
