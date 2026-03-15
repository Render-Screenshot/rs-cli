package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config holds the CLI configuration.
type Config struct {
	APIKey      string `json:"api_key,omitempty"`
	PublicKeyID string `json:"public_key_id,omitempty"`
	SecretKey   string `json:"secret_key,omitempty"`
}

// Path returns the path to the config file.
// Override with RS_CONFIG_PATH environment variable.
func Path() string {
	if p := os.Getenv("RS_CONFIG_PATH"); p != "" {
		return p
	}
	return filepath.Join(configDir(), "renderscreenshot", "config.json")
}

// Load reads the config from disk. Returns empty config if file doesn't exist.
func Load() (*Config, error) {
	path := Path()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

// Save writes the config to disk with restricted permissions.
func Save(cfg *Config) error {
	path := Path()
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// Delete removes the config file.
func Delete() error {
	path := Path()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing config: %w", err)
	}
	return nil
}

// Get returns a config value by key name.
func (c *Config) Get(key string) (string, error) {
	switch key {
	case "api_key":
		return c.APIKey, nil
	case "public_key_id":
		return c.PublicKeyID, nil
	case "secret_key":
		return c.SecretKey, nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}

// Set updates a config value by key name.
func (c *Config) Set(key, value string) error {
	switch key {
	case "api_key":
		c.APIKey = value
	case "public_key_id":
		c.PublicKeyID = value
	case "secret_key":
		c.SecretKey = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return nil
}

// ResolveAPIKey returns the API key using precedence:
// flag > env > config file.
func ResolveAPIKey(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if v := os.Getenv("RS_API_KEY"); v != "" {
		return v
	}
	cfg, err := Load()
	if err != nil {
		return ""
	}
	return cfg.APIKey
}

// ResolveSigningKeys returns the public key ID and secret key.
func ResolveSigningKeys() (publicKeyID, secretKey string) {
	if v := os.Getenv("RS_PUBLIC_KEY_ID"); v != "" {
		publicKeyID = v
	}
	if v := os.Getenv("RS_SECRET_KEY"); v != "" {
		secretKey = v
	}
	if publicKeyID != "" && secretKey != "" {
		return
	}
	cfg, err := Load()
	if err != nil {
		return
	}
	if publicKeyID == "" {
		publicKeyID = cfg.PublicKeyID
	}
	if secretKey == "" {
		secretKey = cfg.SecretKey
	}
	return
}

func configDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("APPDATA")
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}
