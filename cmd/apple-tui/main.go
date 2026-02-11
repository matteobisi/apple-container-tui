package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"container-tui/src/services"
	"container-tui/src/ui"
)

var version = "dev"

func main() {
	var dryRun bool

	rootCmd := &cobra.Command{
		Use:   "apple-tui",
		Short: "Apple Container TUI",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := services.CheckCLI(context.Background()); err != nil {
				return err
			}

			var executor services.CommandExecutor
			if dryRun {
				executor = services.DryRunExecutor{}
			} else {
				executor = services.RealExecutor{}
			}

			configManager, err := services.NewConfigManager()
			if err != nil {
				return err
			}
			config, _, err := configManager.Load()
			if err != nil {
				return err
			}
			ui.ApplyTheme(config.ThemeMode)
			logWriter, err := services.NewLogWriter(config.LogRetentionDays)
			if err != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "warning: failed to initialize command log writer")
			} else {
				executor = services.NewLoggingExecutor(executor, logWriter, dryRun)
			}

			if !dryRun {
				statusBuilder := services.CheckDaemonStatusBuilder{}
				statusCmd, buildErr := statusBuilder.Build()
				if buildErr == nil {
					result, execErr := executor.Execute(statusCmd)
					if execErr == nil {
						status := services.ParseDaemonStatus(result.Stdout)
						if !status.Running {
							_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "Container daemon is not running. Use the Manage screen to start it.")
						}
					}
				}
			}

			program := tea.NewProgram(ui.NewAppModel(executor, version))
			if _, err := program.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview commands without executing")

	rootCmd.Version = version
	rootCmd.SetVersionTemplate("apple-tui version {{.Version}}\n")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
