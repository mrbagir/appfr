package appfr

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func (a *App) ParseConfig(s any) {
	if err := env.Parse(s); err != nil {
		a.logger.Fatalf("failed to parse env: %v", err)
	}
}

func (a *App) readConfig() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		a.logger.Fatalf("failed to load .env file: %v", err)
	}
}
