package main

import (
	"log"

	"github.com/kweonminsung/console-cinema/cmd"
	"github.com/kweonminsung/console-cinema/pkg/cache"
)

func main() {
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
