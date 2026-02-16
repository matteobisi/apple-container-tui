package services

import (
	"os"
	"path/filepath"

	"container-tui/src/models"

	"github.com/spf13/viper"
)

// ConfigManager loads configuration from supported paths.
type ConfigManager struct {
	ReadPaths []string
	WritePath string
}

// NewConfigManager builds a manager with macOS default paths.
func NewConfigManager() (*ConfigManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	readPaths := []string{
		filepath.Join(home, ".config", "actui", "config"),
		filepath.Join(home, "Library", "Application Support", "actui", "config"),
	}
	writePath := filepath.Join(home, "Library", "Application Support", "actui", "config")

	return &ConfigManager{ReadPaths: readPaths, WritePath: writePath}, nil
}

// Load reads the first available config file or returns defaults.
func (m *ConfigManager) Load() (models.UserConfig, string, error) {
	config := models.DefaultUserConfig()
	for _, path := range m.ReadPaths {
		if _, err := os.Stat(path); err != nil {
			continue
		}

		v := viper.New()
		v.SetConfigFile(path)
		v.SetConfigType("toml")
		v.SetDefault("default_build_file", config.DefaultBuildFile)
		v.SetDefault("confirm_destructive_actions", config.ConfirmDestructiveActions)
		v.SetDefault("theme_mode", config.ThemeMode)
		v.SetDefault("refresh_on_focus", config.RefreshOnFocus)
		v.SetDefault("log_retention_days", config.LogRetentionDays)

		if err := v.ReadInConfig(); err != nil {
			return config, path, err
		}
		if err := v.Unmarshal(&config); err != nil {
			return config, path, err
		}
		return config, path, nil
	}

	return config, "", nil
}
