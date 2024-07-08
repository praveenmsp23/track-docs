package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/praveenmsp23/trackdocs/pkg/logger"
)

type ApplicationEnv string

const (
	ApplicationEnvLocal      ApplicationEnv = "local"
	ApplicationEnvProduction ApplicationEnv = "production"
)

var (
	// ErrMissingEnvironmentDatasource missing datasource configuration
	ErrMissingEnvironmentDatasource = errors.New("missing datasource ENV Variable")
)

// Build info string for build info
var (
	BuildSHA    string
	BuildBranch string
	BuildTime   string
)

// Config for the environment
type Config struct {
	Port                 string         `envconfig:"PORT" default:"8080"`
	Listen               string         `envconfig:"LISTEN" default:"0.0.0.0"`
	Env                  ApplicationEnv `envconfig:"ENV" default:"local"`
	APIUrl               string         `envconfig:"API_URL" default:"http://api:8080"`
	TokenHeader          string         `envconfig:"TOKEN_HEADER" default:"X-Access-Token"`
	TokenProvider        string         `envconfig:"TOKEN_PROVIDER" default:"redis"`
	TokenLifeTime        int64          `envconfig:"TOKEN_LIFETIME" default:"86400"`
	CacheSource          string         `envconfig:"CACHE_SOURCE" default:"redis:6379"`
	CacheSourcePassword  string         `envconfig:"CACHE_SOURCE_PASSWORD" default:"password"`
	MeilisearchHost      string         `envconfig:"MEILISEARCH_HOST" default:"http://meiliesearch:7700"`
	MeilisearchMasterKey string         `envconfig:"MEILISEARCH_MASTER_KEY" default:"master_key"`
	Secret               string         `envconfig:"SECRET" default:"k;r(>.]kW6M#NCXK=<EF&}an1JW9!q"` // encrypt and decrypt
	Datasource           string         `envconfig:"DATASOURCE" default:"host=localhost user=trackdocs password=trackdocs dbname=trackdocs port=5432 sslmode=disable TimeZone=Asia/Kolkata"`
}

// NewConfig reads configuration from environment variables and validates it
func NewConfig() (*Config, error) {
	cfg := new(Config)
	err := envconfig.Process("trackdocs", cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse environment variables")
	}
	err = cfg.validate()
	if err != nil {
		return nil, errors.Wrap(err, "config validation failed")
	}

	if cfg.Env == ApplicationEnvLocal {
		logger.LocalInit()
	}
	return cfg, nil
}

func (cfg *Config) validate() error {
	if cfg.Datasource == "" {
		return ErrMissingEnvironmentDatasource
	}
	return nil
}
