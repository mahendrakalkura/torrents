package settings

import (
	"github.com/BurntSushi/toml"
)

type container struct {
	Gorilla containerGorilla `toml:"gorilla"`
	Others  containerOthers  `toml:"others"`
}

type containerGorilla struct {
	Hostname string `toml:"hostname"`
	Port     string `toml:"port"`
}

type containerOthers struct {
	Environment string `toml:"environment"`
	URL         string `toml:"url"`
}

// Container ...
var Container *container

func init() {
	Container = &container{}
	_, err := toml.DecodeFile("settings.toml", Container)
	if err != nil {
		panic(err)
	}
}
