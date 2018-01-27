package main

import (
	"flag"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"github.com/revel/config"
	"github.com/solo-io/glue/cache"
	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/translator"
	"github.com/solo-io/glue/xds"
)

func main() {
	configFile := flag.String("c", "gateway-config.yml", "file to watch for (hot) config")
	port := flag.Int("p", 8081, "xds server port")
	flag.Parse()
	log.Fatalf("%v", start(*configFile, *port))
}

func start() error {
	errChan := make(chan error)
	configChanged := make(chan bool)
	var c cache.Cache
	var t translator.Translator
	go func() {
		errChan <- initTranslator()
	}()
	go func() {
		errChan <- initConfigCache()
	}()
	//TODO: endpoint discovery?
	go func() {
		errChan <- xds.RunXDS(gatewayConfig, xdsPort, configChanged)
	}()
	return <-errChan
}

func watchConfigChanges(gatewayConfig *Config, configFile string, configChanged chan bool) error {
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
