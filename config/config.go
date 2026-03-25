package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

const defaultConfigPath = "~/.passall/config.yaml"

// Config holds the application configuration loaded from the YAML config file.
type Config struct {
	Storage StorageConfig `yaml:"storage"`
	Auth    AuthConfig    `yaml:"auth"`
}

// StorageConfig holds storage-related configuration.
type StorageConfig struct {
	VaultDir string `yaml:"vault_dir"`
}

// AuthConfig holds authentication-related configuration.
type AuthConfig struct {
	MasterPasswordHashFile string `yaml:"master_password_hash_file"`
}

// defaultConfig returns the built-in convention defaults.
// These apply when no config file exists, giving users an out-of-box experience.
func defaultConfig() Config {
	return Config{
		Storage: StorageConfig{
			VaultDir: "~/.passall/vault",
		},
		Auth: AuthConfig{
			MasterPasswordHashFile: "~/.passall/master_password_hash.json",
		},
	}
}

// Load reads and parses the config file at the default path (~/.passall/config.yaml).
// 加载默认配置文件
func Load() (*Config, error) {
	return LoadFrom(defaultConfigPath)
}

// LoadFrom reads and parses the config file at the given path.
// If the file does not exist, the built-in convention defaults are used so users
// get an out-of-box experience without creating a config file manually.
// All "~/" path prefixes in field values are expanded before returning.
func LoadFrom(path string) (*Config, error) {
	//拓展home符，home符号是shell语法糖，操作系统不识别
	expandedPath, err := expandHomePath(path)
	if err != nil {
		return nil, fmt.Errorf("config: resolve path %q: %w", path, err)
	}

	// Start with convention defaults so missing fields are always populated.
	cfg := defaultConfig()

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// No config file — use defaults as-is (convention over configuration).
			if err := cfg.expandPaths(); err != nil {
				return nil, err
			}
			return &cfg, nil
		}
		return nil, fmt.Errorf("config: read file %q: %w", expandedPath, err)
	}

	// Config file exists: overlay user values on top of defaults.
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: invalid configuration: %w", err)
	}

	if err := cfg.expandPaths(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// validate checks that required configuration fields are present.
func (c *Config) validate() error {
	if c.Storage.VaultDir == "" {
		return errors.New("storage.vault_dir must not be empty")
	}
	if c.Auth.MasterPasswordHashFile == "" {
		return errors.New("auth.master_password_hash_file must not be empty")
	}
	return nil
}

// expandPaths expands "~/" prefixes in all path fields in place.
func (c *Config) expandPaths() error {
	var err error
	if c.Storage.VaultDir, err = expandHomePath(c.Storage.VaultDir); err != nil {
		return fmt.Errorf("config: expand vault_dir: %w", err)
	}
	if c.Auth.MasterPasswordHashFile, err = expandHomePath(c.Auth.MasterPasswordHashFile); err != nil {
		return fmt.Errorf("config: expand master_password_hash_file: %w", err)
	}
	return nil
}

// expandHomePath replaces a leading "~/" with the current user's home directory.
func expandHomePath(path string) (string, error) {
	if len(path) < 2 || path[:2] != "~/" {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}
	return filepath.Join(home, path[2:]), nil
}
