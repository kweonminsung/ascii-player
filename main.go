package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kweonminsung/console-cinema/cmd"
	"github.com/kweonminsung/console-cinema/pkg/cache"
)

func main() {
	logDir, err := cache.GetLogDir()
	if err != nil {
		log.Fatalf("failed to get log directory: %v", err)
	}
	logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	cmd.Execute()
}
