package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Items []Item

type Item struct {
	Repo  string `yaml:"repo"`
	Build string `yaml:"build"`
}

func LoadConfig(file string) (res Items, err error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(b, &res)
	return
}
