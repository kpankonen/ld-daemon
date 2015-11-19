package main

import (
	"flag"
	"fmt"
	ld "github.com/launchdarkly/go-client"
	ldr "github.com/launchdarkly/go-client/redis"
	"gopkg.in/gcfg.v1"
	"os"
	"time"
)

var VERSION = "DEV"

type Config struct {
	Redis struct {
		Host string
		Port int
	}
	Main struct {
		ApiKey      string
		Prefix      string
		ExitOnError bool
		StreamUri   string
		BaseUri     string
	}
}

var configFile string

var client ld.LDClient

func main() {
	flag.StringVar(&configFile, "config", "/etc/ld-daemon.conf", "configuration file location")

	flag.Parse()

	var c Config

	fmt.Printf("Starting LaunchDarkly daemon version %s with config %s\n", VERSION, configFile)

	err := gcfg.ReadFileInto(&c, configFile)

	if err != nil {
		fmt.Println("Failed to read configuration file. Exiting.")
		os.Exit(1)
	}

	clientConfig := ld.DefaultConfig
	clientConfig.Stream = true
	clientConfig.FeatureStore = ldr.NewRedisFeatureStore(c.Redis.Host, c.Redis.Port, c.Main.Prefix)
	clientConfig.StreamUri = c.Main.StreamUri
	clientConfig.BaseUri = c.Main.BaseUri

	client = ld.MakeCustomClient(c.Main.ApiKey, clientConfig)
	client.InitializeStream()

	init := make(chan bool)

	go func() {
		for {
			if client.IsStreamInitialized() {
				init <- true
				break
			}
		}
	}()

loop:
	for {
		select {
		case <-init:
			fmt.Println("Initialized stream")
			break loop
		case <-time.After(time.Second * 10):
			fmt.Println("Timed out connecting to stream")
			if c.Main.ExitOnError {
				os.Exit(1)
			}
		}
	}

	for {
		time.Sleep(time.Second)
		if client.IsStreamDisconnected() {
			fmt.Println("Stream connection lost")
			if c.Main.ExitOnError {
				os.Exit(1)
			}
		}
	}

}
