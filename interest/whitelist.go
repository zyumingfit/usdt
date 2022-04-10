package interest

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Whitelists struct {
	Whitelist []Whitelist `yaml:"whitelist"`
}

type Whitelist struct {
	Address string `yaml:"address"`
}

//Get whitelist form config file: whitelist.yaml.
func GetWhitelist() ([]Whitelist, error) {
	whitelists := Whitelists{}

	yamlFile, err := ioutil.ReadFile("./whitelist.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &whitelists)
	if err != nil {
		return nil, err
	}
	return whitelists.Whitelist, nil
}
