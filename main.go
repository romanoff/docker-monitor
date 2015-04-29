package main

import (
	"fmt"
	"os"
)

var config *Config

func main() {
	config, err := ParseConfig("docker-monitor.toml")
	if err != nil {
		fmt.Printf("Error while parsing docker-monitor.toml config: %v\n", err.Error())
		os.Exit(1)
	}
}
