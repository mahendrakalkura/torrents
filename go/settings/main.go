package settings

import (
	"github.com/BurntSushi/toml"
)

type container struct {
	Gorilla containerGorilla `toml:"gorilla"`
	Spiders containerSpiders `toml:"spiders"`
	Others  containerOthers  `toml:"others"`
}

type containerGorilla struct {
	Hostname string `toml:"hostname"`
	Port     string `toml:"port"`
}

type containerSpiders struct {
	URLs []string `toml:"urls"`
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
