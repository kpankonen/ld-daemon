package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	ld "github.com/launchdarkly/go-client"
	ldr "github.com/launchdarkly/go-client/redis"
	"os"
	"time"
)

type Config struct {
	Redis struct {
		Host string
		Port int
	}
	Main struct {
		ApiKey string
		Prefix string
	}
}

var client ld.LDClient

func main() {
	var c Config
	err := gcfg.ReadFileInto(&c, "ldd.conf")

	if err != nil {
		fmt.Println("Failed to read configuration file")
		os.Exit(1)
	}

	clientConfig := ld.DefaultConfig
	clientConfig.Stream = true
	clientConfig.FeatureStore = ldr.NewRedisFeatureStore(c.Redis.Host, c.Redis.Port, c.Main.Prefix)

	client = ld.MakeCustomClient(c.Main.ApiKey, clientConfig)

	init := make(chan bool)

	go func() {
		for {
			if client.IsStreamInitialized() {
				init <- true
				break
			}
		}
	}()

	select {
	case <-init:
		fmt.Println("Initialized stream")
	case <-time.After(time.Second * 10):
		fmt.Println("Timed out connecting to stream")
		os.Exit(1)
	}

	for {
		time.Sleep(time.Second)
		if client.IsStreamDisconnected() {
			fmt.Println("Stream connection lost, exiting")
			os.Exit(1)
		}
	}

}
