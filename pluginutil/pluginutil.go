package pluginutil

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"

	"github.com/fsnotify/fsnotify"
)

// CompilePlugin compiles .go files in a directory into a plugin.
func CompilePlugin(dir string, output string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return err
	}

	args := append([]string{"build", "-buildmode=plugin", "-o", output}, files...)
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// LoadPlugin loads a compiled plugin from a file.
func LoadPlugin(path string) (*plugin.Plugin, error) {
	return plugin.Open(path)
}

// GetSymbols returns a list of symbols defined in the plugin.
func GetSymbols(p *plugin.Plugin) ([]string, error) {
	symGetSymbols, err := p.Lookup("GetSymbols")
	if err != nil {
		return nil, err
	}

	getSymbols := symGetSymbols.(func() []string)
	return getSymbols(), nil
}

// WatchDirectory watches a directory for changes in .go files and recompiles and reloads the plugin.
func WatchDirectory(dir string, output string, onReload func(*plugin.Plugin)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	// err = watcher.Add(dir + "/plugin.go")
	// if err != nil {
	// 	return err
	// }

	// done := make(chan bool)
	// go func() {
	log.Printf("one time listen")
	for event := range watcher.Events {
		log.Printf("Event: %+v", event)
		if event.Op&fsnotify.Write == fsnotify.Write && filepath.Ext(event.Name) == ".go" {
			fmt.Println("File modified:", event.Name)
			if err := CompilePlugin(dir, output); err != nil {
				fmt.Println("Error compiling plugin:", err)
				continue
			}
			newPlugin, err := LoadPlugin(output)
			if err != nil {
				fmt.Println("Error loading plugin:", err)
				continue
			}
			p := newPlugin
			fmt.Println("Plugin reloaded")
			onReload(p)
		}
	}
	// }()

	return nil
}
