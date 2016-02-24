package main

import (
	"flag"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/kelseyhightower/envconfig"
	ld "github.com/launchdarkly/go-client"
	ldr "github.com/launchdarkly/go-client/redis"
	"gopkg.in/gcfg.v1"
	"os"
	"strings"
	"time"
)

var VERSION = "DEV"

const CONFIG_ARG = "config"

type Config struct {
	Redis struct {
		Host string `envconfig:"redis_host"`
		Port int    `envconfig:"redis_port"`
		Uri  string `envconfig:"redis_uri"`
	}
	Main struct {
		ApiKey      string
		Prefix      string
		ExitOnError bool   `envconfig:"exit_on_error"`
		StreamUri   string `envconfig:"stream_uri"`
		BaseUri     string `envconfig:"base_uri"`
		EventsUri   string `envconfig:"events_uri"`
	}
}

var configFile string

var client *ld.LDClient

func exists(file string) bool {
	if _, err := os.Stat("/path/to/whatever"); err == nil {
		return true
	}
	return false
}

func main() {
	flag.StringVar(&configFile, CONFIG_ARG, "/etc/ld-daemon.conf", "configuration file location")

	flag.Parse()

	var envConfig, c Config

	var isConfigSet bool

	flag.Visit(func(f *flag.Flag) {
		if f.Name == CONFIG_ARG {
			isConfigSet = true
		}
	})

	if isConfigSet || exists(configFile) {
		fmt.Printf("Starting LaunchDarkly daemon version %s with configuration file %s\n", formatVersion(VERSION), configFile)
	} else {
		fmt.Printf("Starting LaunchDarkly daemon version %s (no configuration file specified)\n", formatVersion(VERSION))
	}

	err := gcfg.ReadFileInto(&c, configFile)

	if err != nil && isConfigSet {
		fmt.Println("Failed to read configuration file. Exiting.")
		os.Exit(1)
	}

	envErr := envconfig.Process("ld", &envConfig.Main)

	if envErr != nil {
		fmt.Println("Invalid environment variable value: %+v", envErr)
		os.Exit(1)
	}
	envErr = envconfig.Process("ld", &envConfig.Redis)
	if envErr != nil {
		fmt.Println("Invalid environment variable value: %+v", envErr)
		os.Exit(1)
	}

	mergeErr := mergo.MergeWithOverwrite(&c, envConfig)
	if mergeErr != nil {
		fmt.Printf("Failed to merge environment variables into configuration: %+v", mergeErr)
		os.Exit(1)
	}

	clientConfig := ld.DefaultConfig
	clientConfig.Stream = true
	if c.Redis.Uri == "" {
		clientConfig.FeatureStore = ldr.NewRedisFeatureStore(c.Redis.Host, c.Redis.Port, c.Main.Prefix, 0)
	} else {
		clientConfig.FeatureStore = ldr.NewRedisFeatureStoreFromUrl(c.Redis.Uri, c.Main.Prefix, 0)
	}
	if c.Main.StreamUri != "" {
		clientConfig.StreamUri = c.Main.StreamUri
	}
	if c.Main.BaseUri != "" {
		clientConfig.BaseUri = c.Main.BaseUri
	}
	if c.Main.EventsUri != "" {
		clientConfig.EventsUri = c.Main.EventsUri
	}

	client, err = ld.MakeCustomClient(c.Main.ApiKey, clientConfig, 2*time.Minute)

	if err != nil && c.Main.ExitOnError {
		os.Exit(1)
	} else {
		fmt.Println("LaunchDarkly connection initialized")
	}

	done := make(chan bool)
	<-done // Block forever
}

func formatVersion(version string) string {
	split := strings.Split(version, "+")

	if len(split) == 2 {
		return fmt.Sprintf("%s (build %s)", split[0], split[1])
	}
	return version
}
