package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config holds runtime settings for the SIP PBX core.
type Config struct {
	SIPBindHost     string
	SIPAdvertisedHost string
	SIPPort           string
	HTTPPort          string
	BufferSize        int
	LogLevel          string
	ServiceName       string
}

// Load reads optional .env files and environment variables.
func Load() (Config, error) {
	loadDotEnv(".env")
	loadDotEnv(".env.local")

	cfg := Config{
		SIPBindHost:       getEnv("SIP_BIND_HOST", "0.0.0.0"),
		SIPAdvertisedHost: getEnv("SIP_ADVERTISED_HOST", "127.0.0.1"),
		SIPPort:           getEnv("PORT_SERVER_SIP", "5060"),
		HTTPPort:          getEnv("HTTP_PORT", "8080"),
		LogLevel:          strings.ToLower(getEnv("LOG_LEVEL", "info")),
		ServiceName:       getEnv("SERVICE_NAME", "kvoip-pbx"),
	}

	buffer, err := strconv.Atoi(getEnv("SIP_BUFFER_SIZE", "8192"))
	if err != nil || buffer <= 0 {
		return Config{}, fmt.Errorf("SIP_BUFFER_SIZE inválido: %q", getEnv("SIP_BUFFER_SIZE", "8192"))
	}
	cfg.BufferSize = buffer
	return cfg, nil
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf("%s:%s", c.SIPBindHost, c.SIPPort)
}

func (c Config) AdvertisedAddr() string {
	return fmt.Sprintf("%s:%s", c.SIPAdvertisedHost, c.SIPPort)
}

func (c Config) HTTPListenAddr() string {
	return fmt.Sprintf(":%s", c.HTTPPort)
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func loadDotEnv(filename string) {
	path := filename
	if !filepath.IsAbs(path) {
		if cwd, err := os.Getwd(); err == nil {
			path = filepath.Join(cwd, filename)
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}
