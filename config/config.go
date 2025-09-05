package config

import (
	_ "embed"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	envBookDB = "PG_BOOK"
)

type Config struct {
	App      string         `yaml:"app"`
	Stack    string         `yaml:"stack"`
	Server   serverConfig   `yaml:"server"`
	Database databaseConfig `yaml:"database"`
	JWT      jwtConfig      `yaml:"jwt"`
}

type jwtConfig struct {
	SecretKey string        `yaml:"secret_key"`
	TTL       time.Duration `yaml:"ttl"`
}

type serverConfig struct {
	Port              int           `yaml:"port"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
}

type databaseConfig struct {
	BookDB string `yaml:"bookDB"`
}

var (
	//go:embed config.yaml
	file    string
	decoder = yaml.NewDecoder(strings.NewReader(file))
)

func NewConfig() *Config {
	cfg := &Config{}
	if err := decoder.Decode(cfg); err != nil {
		log.Fatal("Config decoding error", err)
	}
	if envVal, exist := os.LookupEnv(envBookDB); exist {
		log.Println("Database variable found")
		cfg.Database.BookDB = envVal
	}
	if err := validate(cfg); err != nil {
		log.Fatal("Wrong configuration", err)
	}
	return cfg
}

func validate(cfg *Config) error {
	switch {
	case cfg.Database.BookDB == "":
		return errors.New("bookDB connection string is empty")
	case cfg.Server.Port == 0:
		return errors.New("server port is zero")
	default:
		return nil
	}
}
