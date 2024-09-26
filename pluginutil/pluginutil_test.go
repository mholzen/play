package pluginutil

import (
	"log"
	"plugin"
	"testing"
)

func Test_Watcher(t *testing.T) {
	t.Skip()
	err := WatchDirectory("/Users/marchome/develop/mholzen/play/plugins", "plugin.so", func(p *plugin.Plugin) {
		symbols, err := GetSymbols(p)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Symbols:", symbols)
	})
	if err != nil {
		log.Fatal(err)
	}
}
