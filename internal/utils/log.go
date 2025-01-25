package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	gap "github.com/muesli/go-app-paths"
)

func getLogFilePath() (string, error) {
	dir, err := gap.NewScope(gap.User, "merlion").CacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "merlion.log"), nil
}

func SetupLog() (func() error, error) {
	log.SetOutput(io.Discard)
	// Log to file, if set
	logFile, err := getLogFilePath()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(logFile), 0o755); err != nil {
		// log disabled
		return func() error { return nil }, nil
	}
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		// log disabled
		return func() error { return nil }, nil
	}
	log.SetOutput(f)
	setLogLevel()
	return f.Close, nil
}

func setLogLevel() {
	level, exists := os.LookupEnv("LOG_LEVEL")
	if !exists {
		log.SetLevel(log.DebugLevel)
		return
	}

	lvl, err := log.ParseLevel(level)
	if err != nil {
		log.SetLevel(log.DebugLevel)
		return
	}

	log.SetLevel(lvl)
}
