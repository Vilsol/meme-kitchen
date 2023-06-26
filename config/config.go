package config

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitializeConfig() context.Context {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("meme")

	initializeDefaults()

	err := viper.ReadInConfig() //nolint:ifshort

	var out io.Writer = os.Stdout
	if !viper.GetBool("production") {
		out = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	log.Logger = zerolog.New(out).With().Str("service", "api").Timestamp().Logger()
	ctx := log.Logger.WithContext(context.Background())

	if err != nil {
		log.Warn().Err(err).Msg("config initialized using defaults and environment only!")
	}

	log.Info().Msg("Config initialized")

	return ctx
}

func initializeDefaults() {
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", "5020")

	viper.SetDefault("production", true)
	viper.SetDefault("profiler", false)

	viper.SetDefault("database.redis.host", "localhost")
	viper.SetDefault("database.redis.port", 6379)
	viper.SetDefault("database.redis.pass", "")
	viper.SetDefault("database.redis.db", 1)

	viper.SetDefault("database.postgres.host", "localhost")
	viper.SetDefault("database.postgres.port", 5432)
	viper.SetDefault("database.postgres.user", "postgres")
	viper.SetDefault("database.postgres.pass", "REPLACE_ME")
	viper.SetDefault("database.postgres.db", "postgres")

	viper.SetDefault("storage.bucket", "smr")
	viper.SetDefault("storage.key", "REPLACE_ME_KEY")
	viper.SetDefault("storage.secret", "REPLACE_ME_SECRET")
	viper.SetDefault("storage.endpoint", "http://localhost:9000")
	viper.SetDefault("storage.region", "eu-central-1")
	viper.SetDefault("storage.base_url", "http://localhost:9000")
}
