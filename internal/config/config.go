package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/zlog"
)

const (
	path    = "./env/config.yaml"
	envPath = "./.env"
)

func Init() *Config {
	wbCfg := config.New()

	err := wbCfg.Load(path, envPath, "")
	if err != nil {
		zlog.Logger.Panic().Err(err).Msg("could not read config file")
	}

	var cfg Config
	if err = wbCfg.Unmarshal(&cfg); err != nil {
		zlog.Logger.Panic().Err(err).Msg("could not unmarshal config file")
	}

	zlog.Logger.Info().Msgf("config: %+v", cfg)

	if err = godotenv.Load(".env"); err != nil {
		zlog.Logger.Warn().Err(err).Msg(".env not found; relying on environment variables")
	}

	val, _ := os.LookupEnv("DB_PASSWORD")
	cfg.Postgres.Password = val

	return &cfg
}
