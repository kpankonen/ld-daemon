package main

import (
	"flag"
	"fmt"
	ld "github.com/launchdarkly/go-client"
	ldr "github.com/launchdarkly/go-client/redis"
	"gopkg.in/gcfg.v1"
	"os"
	"strings"
	"time"
)

var VERSION = "DEV"

type EnvConfig struct {
	ApiKey string
	Prefix string
}

type Config struct {
	Redis struct {
		Host string
		Port int
	}
	Main struct {
		ExitOnError bool
		StreamUri   string
		BaseUri     string
	}
	Environment map[string]*EnvConfig
}

var configFile string

func main() {
	flag.StringVar(&configFile, "config", "/etc/ld-daemon.conf", "configuration file location")

	flag.Parse()

	var c Config

	fmt.Printf("Starting LaunchDarkly daemon version %s with configuration file %s\n", formatVersion(VERSION), configFile)

	err := gcfg.ReadFileInto(&c, configFile)

	if err != nil {
		fmt.Println("Failed to read configuration file. Exiting.")
		os.Exit(1)
	}

	for envName, envConfig := range c.Environment {
		go func(envName string, envConfig EnvConfig) {
			clientConfig := ld.DefaultConfig
			clientConfig.Stream = true
			clientConfig.FeatureStore = ldr.NewRedisFeatureStore(c.Redis.Host, c.Redis.Port, envConfig.Prefix, 10*time.Second)
			clientConfig.StreamUri = c.Main.StreamUri
			clientConfig.BaseUri = c.Main.BaseUri

			_, err := ld.MakeCustomClient(envConfig.ApiKey, clientConfig, time.Second*10)

			if err != nil {
				fmt.Printf("Error initializing LaunchDarkly client for %s: %+v\n", envName, err)

				if c.Main.ExitOnError {
					os.Exit(1)
				}
			} else {
				fmt.Printf("Initialized LaunchDarkly client for %s\n", envName)
			}
		}(envName, *envConfig)
	}

	go forever()
	select {} // block forever
}

func forever() {
	for {
		time.Sleep(time.Second)
	}
}

func formatVersion(version string) string {
	split := strings.Split(version, "+")

	if len(split) == 2 {
		return fmt.Sprintf("%s (build %s)", split[0], split[1])
	}
	return version
}
