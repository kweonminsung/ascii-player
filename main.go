package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/kweonminsung/console-cinema/cmd"
	"github.com/kweonminsung/console-cinema/pkg/cache"
)

func main() {
	cacheDir, err := cache.GetCacheDir()
	if err != nil {
		log.Fatalf("failed to get cache directory: %v", err)
	}
	logPath := filepath.Join(cacheDir, "error.log")
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	// Clear cache at the start of the program
	if err := cache.ClearCache(); err != nil {
		log.Printf("Failed to clear cache at start: %v", err)
	}

	// Defer clearing cache until the program exits
	defer func() {
		if err := cache.ClearCache(); err != nil {
			log.Printf("Failed to clear cache at exit: %v", err)
		}
	}()

	cmd.Execute()
}
