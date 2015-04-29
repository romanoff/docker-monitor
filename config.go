package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type Config struct {
	Repositories map[string]*Repository
	Registries   map[string]*Registry
	Dockerfiles  map[string]*Dockerfile
}

type Repository struct {
	Url    string
	Branch string
}

type Registry struct {
	Host     string
	Port     string
	Protocol string
}

type Dockerfile struct {
	Path         string
	Repositories []string
	Registries   []string
}

func ParseConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	_, err = toml.Decode(string(content), &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
