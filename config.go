package main

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Interval     string
	Repositories map[string]*Repository
	Registries   map[string]*Registry
	Dockerfiles  map[string]*Dockerfile
	Exit         chan bool
}

func (self *Config) startCron() {
	IntervalString := strings.Trim(self.Interval, "m")
	interval, err := strconv.Atoi(IntervalString)
	if err != nil {
		interval = 5
	}
	for _ = range time.Tick(time.Duration(interval) * time.Minute) {
		log.Println("Checking repositories")
		self.CheckRepositories()
	}
	// To exit
	// self.Exit <- true
}

func (self *Config) CheckRepositories() {
	var wg sync.WaitGroup
	for _, repository := range self.Repositories {
		wg.Add(1)
		go func(repository *Repository) {
			defer wg.Done()
			err := repository.GetLatestSha()
			if err != nil {
				log.Fatal(err)
			}
		}(repository)
	}
	wg.Wait()
	for name, dockerfile := range self.Dockerfiles {
		wg.Add(1)
		go func(dockerfile *Dockerfile) {
			defer wg.Done()
			err := dockerfile.CheckIfUpdated(name)
			if err != nil {
				log.Fatal(err)
			}
		}(dockerfile)
	}
	wg.Wait()
}

type Repository struct {
	Url    string
	Branch string
	Sha    string
}

func (self *Repository) GetLatestSha() error {
	branch := self.Branch
	if branch == "" {
		branch = "master"
	}
	cmd := exec.Command("git", "ls-remote", self.Url, branch)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return err
	}
	if self.Sha != cmdOutput.String()[:40] {
		log.Println(self.Url + " repository updated to '" + self.Sha + "' from '" + self.Sha + "'")
	}
	self.Sha = cmdOutput.String()[:40]
	return nil
}

type Registry struct {
	Host string
	Port string
}

type Dockerfile struct {
	Path            string
	Repositories    []string
	Delay           string
	Registries      []string
	RepositoriesSha string
}

var mutex = &sync.Mutex{}

func (self *Dockerfile) CheckIfUpdated(name string) error {
	repositoriesSha := ""
	for _, repositoryName := range self.Repositories {
		repositoriesSha += config.Repositories[repositoryName].Sha
	}
	if repositoriesSha != self.RepositoriesSha {
		log.Println(name + " dockerfile updated. Rebuilding")
		self.RepositoriesSha = repositoriesSha
		self.Rebuild(name)
		self.PushToRegistries(name)
		err := WriteDockerfileSha(name, repositoriesSha)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Dockerfile) Rebuild(name string) error {
	delayString := strings.Trim(self.Delay, "m")
	delay, err := strconv.Atoi(delayString)
	if err != nil {
		delay = 0
	}
	time.Sleep(time.Duration(delay) * time.Minute)
	cmd := exec.Command("sudo", "docker", "build", "-f", self.Path, "--no-cache", "true", "--force-rm", "true", "-t", name)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (self *Dockerfile) PushToRegistries(name string) error {
	for _, registryName := range self.Registries {
		registry := config.Registries[registryName]
		if registry == nil {
			continue
		}
		host := registry.Host
		port := registry.Port
		if port == "" {
			port = "5000"
		}
		remoteImageName := host + ":" + port + "/" + name
		cmd := exec.Command("sudo", "docker", "tag", name, remoteImageName)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(remoteImageName + " image is pushed to remote repository")
		cmd = exec.Command("sudo", "docker", "push", remoteImageName)
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
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
	conf.Exit = make(chan bool)
	return conf, nil
}
