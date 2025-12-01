package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zarazaex69/zik/apps/cli/internal/config"
	"gopkg.in/yaml.v3"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage ZIK configuration",
		Long:  `View and modify ZIK configuration settings.`,
	}

	configListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE:  runConfigList,
	}

	configGetCmd = &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE:  runConfigGet,
	}

	configSetCmd = &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE:  runConfigSet,
	}
)

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}

func runConfigList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	fmt.Println("Current configuration:")
	fmt.Println("─────────────────────")
	fmt.Print(string(data))

	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	// TODO: Implement config get
	return fmt.Errorf("not implemented yet")
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	// TODO: Implement config set
	return fmt.Errorf("not implemented yet")
}
