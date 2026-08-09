package config

import (
	"errors"
	"os"
)

const (
	defaultDataDir = "data"
	dataDirEnv     = "DATA_DIR"
	tokenEnv       = "DISCORD_TOKEN"
	guildIDEnv     = "GUILD_ID"
)

type Config struct {
	Token   string
	GuildID string
	DataDir string
}

func Load() (*Config, error) {
	token := os.Getenv(tokenEnv)
	if token == "" {
		return nil, errors.New("DISCORD_TOKEN is not set")
	}

	cfg := &Config{
		Token:   token,
		GuildID: os.Getenv(guildIDEnv),
		DataDir: os.Getenv(dataDirEnv),
	}
	if cfg.DataDir == "" {
		cfg.DataDir = defaultDataDir
	}
	return cfg, nil
}
