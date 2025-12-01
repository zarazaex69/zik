package main

import (
	"github.com/spf13/cobra"
)

var (
	chatCmd = &cobra.Command{
		Use:   "chat",
		Short: "Start an interactive chat session with AI",
		Long: `Start an interactive chat session with AI.
Maintains conversation context and allows multi-turn conversations.`,
		RunE: runChat,
	}
)

func runChat(cmd *cobra.Command, args []string) error {
	// TODO: Implement interactive chat with bubbletea
	return nil
}
