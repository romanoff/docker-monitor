package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

var config *Config

func main() {
	config, err := ParseConfig("docker-monitor.toml")
	if err != nil {
		fmt.Printf("Error while parsing docker-monitor.toml config: %v\n", err.Error())
		os.Exit(1)
	}
	db, err = bolt.Open("docker-monitor.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	for name, dockerfile := range config.Dockerfiles {
		sha, err := ReadDockerfileSha(name)
		if err != nil {
			dockerfile.RepositoriesSha = sha
		}
	}
	go config.startCron()
	<-config.Exit
}
