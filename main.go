package main

import (
	"flag"
	"io/ioutil"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/solo-io/glue/config"
	"github.com/solo-io/glue/module"
	_ "github.com/solo-io/glue/module/install"
	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/xds"
)

func main() {
	configFile := flag.String("c", "gateway-config.yml", "file to watch for (hot) config")
	port := flag.Int("p", 8081, "xds server port")
	flag.Parse()
	log.Fatalf("%v", start(*configFile, *port))
}

func start(configFile string, xdsPort int) error {
	errChan := make(chan error)
	configChanged := make(chan bool)
	gatewayConfig := config.NewConfig()
	module.Init(gatewayConfig)
	// give modules a sec to register
	time.Sleep(time.Millisecond * 250)
	go func() {
		errChan <- watchConfigChanges(gatewayConfig, configFile, configChanged)
	}()
	go func() {
		errChan <- xds.RunXDS(gatewayConfig, xdsPort, configChanged)
	}()
	return <-errChan
}

func watchConfigChanges(gatewayConfig *config.Config, configFile string, configChanged chan bool) error {
	if err := setConfig(gatewayConfig, configFile, configChanged); err != nil {
		return err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	if err := watcher.Add(configFile); err != nil {
		return err
	}
	for {
		select {
		case event := <-watcher.Events:
			log.Printf("config changed: %v", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("modified file: %v", event.Name)
			}

			if err := watcher.Add(configFile); err != nil {
				return err
			}
			if err := setConfig(gatewayConfig, configFile, configChanged); err != nil {
				return err
			}
		case err := <-watcher.Errors:
			log.Printf("watcher error: %v", err)
		}
	}
}

func setConfig(gatewayConfig *config.Config, configFile string, configChanged chan bool) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := gatewayConfig.Update(data); err != nil {
		log.GreyPrintf("Warning: config was rejected: \n%s\n with err: %v", data, err)
		return nil
	}

	go func() {
		configChanged <- true
	}()

	log.GreyPrintf("config set:\n%s", string(data))
	return nil
}
