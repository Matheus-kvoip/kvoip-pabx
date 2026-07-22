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
	SIPBindHost       string
	SIPAdvertisedHost string
	SIPPort           string
	HTTPPort          string
	BufferSize        int
	LogLevel          string
	ServiceName       string
	AuthEnabled       bool
	AuthRealm         string
	SIPUsers          map[string]string
	MediaEnabled      bool
	MediaBindHost     string
	MediaAdvertiseHost string
	RTPPortMin        int
	RTPPortMax        int
}

// Load reads optional .env files and environment variables.
func Load() (Config, error) {
	loadDotEnv(".env")
	loadDotEnv(".env.local")

	cfg := Config{
		SIPBindHost:        getEnv("SIP_BIND_HOST", "0.0.0.0"),
		SIPAdvertisedHost:  getEnv("SIP_ADVERTISED_HOST", "127.0.0.1"),
		SIPPort:            getEnv("PORT_SERVER_SIP", "5060"),
		HTTPPort:           getEnv("HTTP_PORT", "8080"),
		LogLevel:           strings.ToLower(getEnv("LOG_LEVEL", "info")),
		ServiceName:        getEnv("SERVICE_NAME", "kvoip-pbx"),
		AuthEnabled:        parseBool(getEnv("SIP_AUTH_ENABLED", "true"), true),
		AuthRealm:          getEnv("SIP_AUTH_REALM", "kvoip.local"),
		SIPUsers:           parseUsers(getEnv("SIP_USERS", "1001:kvoip123,1002:kvoip123")),
		MediaEnabled:       parseBool(getEnv("MEDIA_ENABLED", "true"), true),
		MediaBindHost:      getEnv("MEDIA_BIND_HOST", "0.0.0.0"),
		MediaAdvertiseHost: getEnv("MEDIA_ADVERTISE_HOST", getEnv("SIP_ADVERTISED_HOST", "127.0.0.1")),
	}

	buffer, err := strconv.Atoi(getEnv("SIP_BUFFER_SIZE", "8192"))
	if err != nil || buffer <= 0 {
		return Config{}, fmt.Errorf("SIP_BUFFER_SIZE inválido: %q", getEnv("SIP_BUFFER_SIZE", "8192"))
	}
	cfg.BufferSize = buffer

	rtpMin, err := strconv.Atoi(getEnv("RTP_PORT_MIN", "10000"))
	if err != nil || rtpMin <= 0 {
		return Config{}, fmt.Errorf("RTP_PORT_MIN inválido")
	}
	rtpMax, err := strconv.Atoi(getEnv("RTP_PORT_MAX", "20000"))
	if err != nil || rtpMax <= rtpMin {
		return Config{}, fmt.Errorf("RTP_PORT_MAX inválido")
	}
	cfg.RTPPortMin = rtpMin
	cfg.RTPPortMax = rtpMax
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

func parseBool(raw string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func parseUsers(raw string) map[string]string {
	out := map[string]string{}
	for _, item := range strings.Split(raw, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		user, pass, ok := strings.Cut(item, ":")
		if !ok {
			continue
		}
		user = strings.TrimSpace(user)
		pass = strings.TrimSpace(pass)
		if user == "" {
			continue
		}
		out[user] = pass
	}
	return out
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
