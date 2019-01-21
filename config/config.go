package config

import (
	"fmt"
	"strings"

	"gopkg.in/caarlos0/env.v2"
)

const (
	MIN_TOKEN_LENGTH = 16
)

type Config struct {
	ListenAddress string   `env:"GO_MICROSERVICE_1_LISTEN_ADDRESS" envDefault:":80"`
	HealthFreqSec int      `env:"GO_MICROSERVICE_1_HEALTH_FREQ_SEC" envDefault:"60"`
	EnvName       string   `env:"GO_MICROSERVICE_1_ENV_NAME" envDefault:"dev"`
	Tokens        []string `env:"GO_MICROSERVICE_1_TOKENS"`
	ServiceName   string   `env:"GO_MICROSERVICE_1_SERVICE_NAME" envDefault:"go-microservice-1"`

	MongoDBName           string   `env:"GO_MICROSERVICE_1_MONGO_DB_NAME"`
	MongoDBHosts          []string `env:"GO_MICROSERVICE_1_MONGO_DB_HOSTS"` // ports included here
	MongoDBUser           string   `env:"GO_MICROSERVICE_1_MONGO_DB_USER"`
	MongoDBPassword       string   `env:"GO_MICROSERVICE_1_MONGO_DB_PASS"`
	MongoDBReplicaSet     string   `env:"GO_MICROSERVICE_1_MONGO_DB_REPLICA_SET"`
	MongoDBSource         string   `env:"GO_MICROSERVICE_1_MONGO_DB_AUTH_SOURCE"`
	MongoDBConnUseSSL     bool     `env:"GO_MICROSERVICE_1_MONGO_DB_USE_SSL" envDefault:"true"`
	MongoDBConnTimeoutSec int      `env:"GO_MICROSERVICE_1_MONGO_DB_TIMEOUT_SEC" envDefault:"30"`

	FooAPIHost         string `env:"GO_MICROSERVICE_1_FOO_API_HOST"`

	StatsDAddress string  `env:"GO_MICROSERVICE_1_STATSD_ADDRESS" envDefault:"localhost:8125"`
	StatsDPrefix  string  `env:"GO_MICROSERVICE_1_STATSD_PREFIX" envDefault:"statsd.go-microservice-1.dev"`
	StatsDRate    float32 `env:"GO_MICROSERVICE_1_STATSD_RATE" envDefault:"1.0"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) LoadEnvVars() error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("Unable to fetch env vars: %v", err.Error())
	}

	var errorList []string

	validations := []validation{
		nonEmptyString{c.MongoDBName, "GO_MICROSERVICE_1_MONGO_DB_NAME"},
		nonEmptyStringSlice{c.MongoDBHosts, "GO_MICROSERVICE_1_MONGO_DB_HOSTS"},
		nonEmptyString{s: c.FooAPIHost, name: "GO_MICROSERVICE_1_FOO_API_HOST"},
		nonEmptyStringSlice{s: c.Tokens, name: "GO_MICROSERVICE_1_TOKENS"},
		tokenLength{s: c.Tokens, name: "GO_MICROSERVICE_1_TOKENS"},
	}

	for _, v := range validations {
		if ok, e := v.validate(); !ok {
			errorList = append(errorList, e...)
		}
	}

	if len(errorList) != 0 {
		return fmt.Errorf(strings.Join(errorList, "; "))
	}

	return nil
}

type validation interface {
	validate() (bool, []string)
}

type nonEmptyString struct {
	s, name string
}

func (v nonEmptyString) validate() (bool, []string) {
	if v.s == "" {
		return false, []string{fmt.Sprintf("missing '%s' env var", v.name)}
	}

	return true, nil
}

type nonEmptyStringSlice struct {
	s    []string
	name string
}

func (v nonEmptyStringSlice) validate() (bool, []string) {
	if len(v.s) < 1 {
		return false, []string{fmt.Sprintf("missing '%s' env var", v.name)}
	}

	return true, nil
}

type tokenLength struct {
	s    []string
	name string
}

func (t tokenLength) validate() (bool, []string) {
	var errorList []string

	for _, token := range t.s {
		if len(token) < MIN_TOKEN_LENGTH {
			errorList = append(errorList, fmt.Sprintf("%v token must be at least %v chars long", token, MIN_TOKEN_LENGTH))
		}
	}

	if len(errorList) > 0 {
		return false, errorList
	}

	return true, nil
}
