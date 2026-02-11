package models

// UserConfig stores persisted user preferences.
type UserConfig struct {
	DefaultBuildFile          string `mapstructure:"default_build_file" toml:"default_build_file"`
	ConfirmDestructiveActions bool   `mapstructure:"confirm_destructive_actions" toml:"confirm_destructive_actions"`
	ThemeMode                 string `mapstructure:"theme_mode" toml:"theme_mode"`
	RefreshOnFocus            bool   `mapstructure:"refresh_on_focus" toml:"refresh_on_focus"`
	LogRetentionDays          int    `mapstructure:"log_retention_days" toml:"log_retention_days"`
}

// DefaultUserConfig returns app defaults.
func DefaultUserConfig() UserConfig {
	return UserConfig{
		DefaultBuildFile:          "Containerfile",
		ConfirmDestructiveActions: true,
		ThemeMode:                 "auto",
		RefreshOnFocus:            false,
		LogRetentionDays:          7,
	}
}
