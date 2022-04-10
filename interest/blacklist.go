package interest

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Blacklists struct {
	Blacklist []Blacklist `yaml:"blacklist"`
}

type Blacklist struct {
	Address string `yaml:"address"`
}

//Get blacklist form config file: blacklist.yaml.
func GetBlacklist() ([]Blacklist, error) {
	blacklists := Blacklists{}

	yamlFile, err := ioutil.ReadFile("./blacklist.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &blacklists)
	if err != nil {
		return nil, err
	}
	return blacklists.Blacklist, nil
}
