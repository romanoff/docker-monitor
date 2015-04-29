package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os/exec"
	"bytes"
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

func (self *Repository) GetLatestSha() (string ,error) {
	branch := self.Branch
	if branch == "" {
		branch = "master"
	}
	cmd := exec.Command("git", "ls-remote", self.Url, branch)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return cmdOutput.String()[:40], nil
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
