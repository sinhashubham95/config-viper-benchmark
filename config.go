package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"sync"
)

const (
	ConfigDirectoryEnv = "CONFIG_DIRECTORY"
	MetaFile           = "package"
	JSON               = "json"
	Properties         = "properties"
)

var (
	ConfigDirectory string
	metaViper       *viper.Viper
	configMu        *sync.RWMutex
	configViper     map[string]*viper.Viper
)

func init() {
	ConfigDirectory = os.Getenv(ConfigDirectoryEnv)

	viper.Debug()

	configMu = &sync.RWMutex{}
	configMu.Lock()
	defer configMu.Unlock()
	configViper = make(map[string]*viper.Viper)
}

func initConfigs() {
	metaViper = viper.New()
	metaViper.SetConfigName(MetaFile)
	metaViper.SetConfigType(JSON)
	metaViper.AddConfigPath(ConfigDirectory)
	metaViper.WatchConfig()

	metaViper.OnConfigChange(onMetaFileChanged)

	// load config files
	checkConfigReload()
}

func onMetaFileChanged(e fsnotify.Event) {
	// reload all files
	checkConfigReload()
}

func checkConfigReload() {
	err := metaViper.ReadInConfig()
	if err != nil {
		_ = fmt.Errorf("error reading meta config file %+v", err)
		return
	}
	reloadConfigs()
}

func reloadConfigs() {
	mark := make(map[string]bool)
	props := metaViper.GetStringSlice(Properties)
	if props != nil {
		for _, prop := range props {
			mark[prop] = true
			configMu.RLock()
			_, ok := configViper[prop]
			configMu.RUnlock()
			if !ok {
				v := viper.New()
				v.SetConfigName(prop)
				v.SetConfigType(JSON)
				v.AddConfigPath(ConfigDirectory)
				v.WatchConfig()
				configMu.Lock()
				configViper[prop] = v
				configMu.Unlock()
			}
		}
	}
	configMu.Lock()
	defer configMu.Unlock()
	for prop := range configViper {
		if !mark[prop] {
			delete(configViper, prop)
			return
		}
		err := configViper[prop].ReadInConfig()
		if err != nil {
			_ = fmt.Errorf("error reading config file %s %+v", prop, err)
		}
	}
}
