package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type KatalogAgentConfig struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"error"`

	ClusterName           string `env:"CLUSTER_NAME" envDefault:"unknown"`
	CronSchedule          string `env:"CRON_SCHEDULE" envDefault:"@every 30m"`
	BlacklistedNamespaces string `env:"BLACKLISTED_NAMESPACES" envDefault:"kube-system,flux-system"`

	MetricsEnabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	MetricsPort    int  `env:"METRICS_PORT" envDefault:"8081"`

	KatalogHost            string `env:"KATALOG_HOST" envDefault:"http://localhost:8080"`
	KatalogCertificateFile string `env:"KATALOG_CERTIFICATE_FILE"`
	KatalogKeyFile         string `env:"KATALOG_KEY_FILE"`
	KatalogCAFile          string `env:"KATALOG_CA_FILE"`

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
