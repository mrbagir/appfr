package config

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	defaultFileName = "/.env"
	// defaultOverrideFileName = "/.local.env"
)

func NewEnvFile(configFolder string) Config {
	conf := &EnvLoader{}
	conf.read(configFolder)

	return conf
}

type EnvLoader struct{}

func (e *EnvLoader) read(folder string) {
	var (
		defaultFile = folder + defaultFileName
		// env         = e.Get("APP_ENV")
	)

	// Only Capture initial system environment before any file loading
	initialEnv := e.captureInitialEnv()

	err := godotenv.Load()
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("Failed to load config from file: %v, Err: %v", defaultFile, err)
		}

		log.Printf("Failed to load config from file: %v, Err: %v", defaultFile, err)
	} else {
		log.Printf("Loaded config from file: %v", defaultFile)
	}

	// Reload system environment variables to ensure they take precedence over values loaded from files.
	// This is required because Overload replaces the original system variables, which we need to restore.
	for key, envVar := range initialEnv {
		os.Setenv(key, envVar)
	}
}

// captureInitialEnv captures the initial system environment keys.
func (*EnvLoader) captureInitialEnv() map[string]string {
	initialEnv := make(map[string]string)

	for _, envVar := range os.Environ() {
		key, value, found := strings.Cut(envVar, "=")
		if found {
			initialEnv[key] = value
		}
	}

	return initialEnv
}

func (*EnvLoader) Get(key string) string {
	return os.Getenv(key)
}

func (*EnvLoader) GetOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}
