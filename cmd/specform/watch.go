package main

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/specform/go-sdk/internal"
)

func WatchFiles(files []string, outputDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to initialize watcher: %w", err)
	}
	defer watcher.Close()

	// Add all files to watcher
	for _, file := range files {
		if err := watcher.Add(file); err != nil {
			return fmt.Errorf("failed to watch %s: %w", file, err)
		}
		fmt.Printf("👀 Watching %s...\n", file)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				fmt.Printf("🔄 Change detected: %s\n", event.Name)
				time.Sleep(100 * time.Millisecond) // debounce
				if out, err := internal.CompileSpecFile(event.Name, outputDir); err != nil {
					fmt.Printf("❌ Error recompiling %s: %v\n", event.Name, err)
				} else {
					fmt.Printf("✅ Recompiled %s → %s\n", event.Name, out)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Printf("⚠️ Watch error: %v\n", err)
		}
	}
}
