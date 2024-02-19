// Модуль конфигурации приложения
package config

import (
	"bytes"
	"dario.cat/mergo"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"os"
)

// ServerConfig описывает структуру конфигурации приложения
type ServerConfig struct {
	RunAddr     string `json:"server_address" env:"SERVER_ADDRESS"`
	DatabaseDSN string `json:"database_dsn" env:"DATABASE_DSN"`
	Config      string `json:"-" env:"CONFIG"`
	TLSCertPath string `json:"tls_cert_path" env:"TLS_CERT_PATH"`
	TLSKeyPath  string `json:"tls_key_path" env:"TLS_KEY_PATH"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"debug"`
	EnableHTTPS bool   `json:"enable_https" env:"ENABLE_HTTPS"`
}

var serverConfig ServerConfig

// ServerConfig парсит значения из переменных окружения
func ParseFlags() (*ServerConfig, error) {
	flag.BoolVar(&serverConfig.EnableHTTPS, "s", true, "enable https")
	flag.StringVar(&serverConfig.DatabaseDSN, "d", "", "Data Source Name (DSN)")
	flag.StringVar(&serverConfig.Config, "c", "", "Config json file path")
	flag.StringVar(&serverConfig.TLSCertPath, "l", "./certs/cert.pem", "path to tls cert file")
	flag.StringVar(&serverConfig.TLSKeyPath, "k", "./certs/private.pem", "path to tls key file")
	flag.StringVar(&serverConfig.LogLevel, "g", "", "log level")
	flag.Parse()

	if err := env.Parse(&serverConfig); err != nil {
		return nil, fmt.Errorf("error parsing env variables: %w", err)
	}

	if serverConfig.Config != "" {
		data, err := os.ReadFile(serverConfig.Config)
		if err != nil {
			return nil, fmt.Errorf("error opening config file: %w", err)
		}

		var configFromFile ServerConfig
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(&configFromFile); err != nil {
			return nil, fmt.Errorf("error parsing json file config: %w", err)
		}

		if err := mergo.Merge(&serverConfig, configFromFile); err != nil {
			return nil, fmt.Errorf("cannot merge configs: %w", err)
		}
	}

	return &serverConfig, nil
}
