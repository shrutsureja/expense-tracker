package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
	Backup   BackupConfig   `yaml:"backup"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type AuthConfig struct {
	JWTSecret  string          `yaml:"jwt_secret"`
	TokenExpiry string         `yaml:"token_expiry"`
	SuperAdmin SuperAdminConfig `yaml:"super_admin"`
}

type SuperAdminConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type BackupConfig struct {
	Enabled        bool              `yaml:"enabled"`
	Interval       string            `yaml:"interval"`
	LocalRetention string            `yaml:"local_retention"`
	GoogleDrive    GoogleDriveConfig `yaml:"google_drive"`
}

type GoogleDriveConfig struct {
	CredentialsFile string `yaml:"credentials_file"`
	FolderID        string `yaml:"folder_id"`
}

func Load() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	cfg := &Config{
		Server: ServerConfig{Port: 8080},
		Database: DatabaseConfig{Path: "./data/expense-tracker.db"},
		Auth: AuthConfig{
			TokenExpiry: "360h",
			SuperAdmin: SuperAdminConfig{
				Username: "admin",
			},
		},
		Backup: BackupConfig{
			Interval:       "168h",
			LocalRetention: "1440h",
		},
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.Auth.JWTSecret == "" {
		return nil, fmt.Errorf("auth.jwt_secret is required")
	}
	if cfg.Auth.SuperAdmin.Password == "" {
		return nil, fmt.Errorf("auth.super_admin.password is required")
	}

	return cfg, nil
}

func (c *Config) TokenExpiryDuration() time.Duration {
	d, err := time.ParseDuration(c.Auth.TokenExpiry)
	if err != nil {
		return 360 * time.Hour
	}
	return d
}

func (c *Config) BackupIntervalDuration() time.Duration {
	d, err := time.ParseDuration(c.Backup.Interval)
	if err != nil {
		return 168 * time.Hour
	}
	return d
}

func (c *Config) BackupRetentionDuration() time.Duration {
	d, err := time.ParseDuration(c.Backup.LocalRetention)
	if err != nil {
		return 1440 * time.Hour
	}
	return d
}
