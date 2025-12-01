package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zarazaex69/zik/apps/cli/internal/config"
)

var (
	// Root command
	rootCmd = &cobra.Command{
		Use:   "zik",
		Short: "ZIK - AI Tools for Developers",
		Long: `ZIK is a powerful CLI tool that brings AI assistance directly to your terminal.
Generate commit messages, chat with AI, review code, and more.`,
		Version: config.Version,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(askCmd)
	rootCmd.AddCommand(codeCmd)
	rootCmd.AddCommand(configCmd)
}
