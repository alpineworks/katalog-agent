package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type KatalogAgentConfig struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"error"`

	CronSchedule          string `env:"CRON_SCHEDULE" envDefault:"@every 30m"`
	BlacklistedNamespaces string `env:"BLACKLISTED_NAMESPACES" envDefault:"kube-system,flux-system"`

	MetricsEnabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	MetricsPort    int  `env:"METRICS_PORT" envDefault:"8081"`

	TracingEnabled    bool    `env:"TRACING_ENABLED" envDefault:"false"`
	TracingSampleRate float64 `env:"TRACING_SAMPLERATE" envDefault:"0.01"`
	TracingService    string  `env:"TRACING_SERVICE" envDefault:"katalog-agent"`
	TracingVersion    string  `env:"TRACING_VERSION"`
}

func NewConfig() (*KatalogAgentConfig, error) {
	var cfg KatalogAgentConfig

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
