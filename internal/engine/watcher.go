package engine

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var ignoreDirs = []string{
	".git", ".snapver", "node_modules", "dist", "build", "out", "venv", "__pycache__", ".idea", ".vscode", ".DS_Store",
}

func isPathIgnored(path string) bool {
	for _, ignore := range ignoreDirs {
		if strings.Contains(path, ignore) {
			return true
		}
	}
	return false
}

// watchLoop continuously watches the file system for changes
func (d *Engine) watchLoop() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("❌ fsnotify.NewWatcher error:", err)
	}
	defer watcher.Close()

	err = filepath.Walk(d.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && !isPathIgnored(path) {
			err := watcher.Add(path)
			if err != nil {
				log.Println("❗ Failed to watch:", path, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal("Failed to walk root:", err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				if isPathIgnored(event.Name) {
					continue
				}
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					// Add new directory to watcher
					if !isPathIgnored(event.Name) {
						watcher.Add(event.Name)
					}
				} else {
					d.handleFileChange(event.Name)
				}
			}
		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
	}
}
