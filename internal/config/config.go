package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"error"`

	MetricsEnabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	MetricsPort    int  `env:"METRICS_PORT" envDefault:"8081"`

	Local bool `env:"LOCAL" envDefault:"false"`

	DailyForecastCronSchedule string `env:"DAILY_FORECAST_CRON_SCHEDULE" envDefault:"30 6 * * *"`

	LFPWeatherForecastInferenceAPIURL     string        `env:"LFP_WEATHER_FORECAST_INFERENCE_API_URL"`
	LFPWeatherForecastInferenceAPIAPIKey  string        `env:"LFP_WEATHER_FORECAST_INFERENCE_API_API_KEY"`
	LFPWeatherForecastInferenceAPITimeout time.Duration `env:"LFP_WEATHER_FORECAST_INFERENCE_API_TIMEOUT" envDefault:"5s"`

	BlueskyHost                                  string        `env:"BLUESKY_HOST" envDefault:"https://bsky.social"`
	BlueskyUsername                              string        `env:"BLUESKY_USERNAME"`
	BlueskyPassword                              string        `env:"BLUESKY_PASSWORD"`
	BlueskyAPITimeout                            time.Duration `env:"BLUESKY_API_TIMEOUT" envDefault:"5s"`
	BlueskyTokenRefreshCronSchedule              string        `env:"BLUESKY_TOKEN_REFRESH_CRON_SCHEDULE" envDefault:"*/10 * * * *"`
	BlueskyTokenRefreshForceRefreshWhenExpiresIn time.Duration `env:"BLUESKY_TOKEN_REFRESH_FORCE_REFRESH_WHEN_EXPIRES_IN" envDefault:"10m"`

	TracingEnabled    bool    `env:"TRACING_ENABLED" envDefault:"false"`
	TracingSampleRate float64 `env:"TRACING_SAMPLERATE" envDefault:"0.01"`
	TracingService    string  `env:"TRACING_SERVICE" envDefault:"katalog-agent"`
	TracingVersion    string  `env:"TRACING_VERSION"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
