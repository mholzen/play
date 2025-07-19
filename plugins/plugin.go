package main

import (
	"log"
	"plugin"

	"github.com/mholzen/play/pluginutil"
)

func SetupWatchers() error {
	dir := "/Users/marchome/develop/mholzen/play/plugins"
	var err error
	// TODO: does not return
	err = pluginutil.WatchDirectory(dir, "plugin.so", func(p *plugin.Plugin) {
		s, err2 := pluginutil.GetSymbols(p)
		if err2 != nil {
			log.Printf("Error getting symbols: %v", err2)
			err = err2
			return
		}
		log.Printf("Symbols: %v", s)
	})
	if err != nil {
		log.Printf("Error watching directory: %v", err)
	} else {
		log.Printf("Watching directory: %s", dir)
	}
	return err
}

func GetSymbols() []string {
	return []string{
		"Transition",
		"Delay",
		"RepeatEvery",
		"home",
		"universe",
		"setup",
		"gold",
		"rainbow",
		"moveTomshine",
		"beatDown",
		"moveDownTomshine",
		"twoColors",
		"clock",
		"main",
		"REFRESH",
	}
}
