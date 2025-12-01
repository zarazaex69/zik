package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zarazaex69/zik/apps/cli/internal/ai"
	"github.com/zarazaex69/zik/apps/cli/internal/config"
)

var (
	askStream bool

	askCmd = &cobra.Command{
		Use:   "ask <question>",
		Short: "Ask a quick question to AI",
		Long: `Ask a one-off question to the AI without maintaining conversation context.
Perfect for quick queries, code explanations, or getting instant answers.`,
		Example: `  zik ask "What is the difference between let and const?"
  zik ask "How do I reverse a string in Go?"
  zik ask --stream "Explain async/await in JavaScript"`,
		Args: cobra.MinimumNArgs(1),
		RunE: runAsk,
	}
)

func init() {
	askCmd.Flags().BoolVarP(&askStream, "stream", "s", true, "Stream the response in real-time")
}

func runAsk(cmd *cobra.Command, args []string) error {
	question := strings.Join(args, " ")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize AI client
	aiClient := ai.NewClient(cfg)
	ctx := context.Background()

	messages := []ai.Message{
		{Role: "user", Content: question},
	}

	if askStream {
		// Streaming response
		chunkChan, errChan := aiClient.ChatStream(ctx, messages, cfg.Temperature, cfg.MaxTokens)

		for {
			select {
			case chunk, ok := <-chunkChan:
				if !ok {
					fmt.Println() // New line at end
					return nil
				}
				if chunk.Content != "" {
					fmt.Print(chunk.Content)
				}
			case err := <-errChan:
				if err != nil {
					return fmt.Errorf("AI request failed: %w", err)
				}
			}
		}
	} else {
		// Non-streaming response
		resp, err := aiClient.Chat(ctx, messages, cfg.Temperature, cfg.MaxTokens)
		if err != nil {
			return fmt.Errorf("AI request failed: %w", err)
		}

		if len(resp.Choices) == 0 {
			return fmt.Errorf("no response from AI")
		}

		fmt.Println(resp.Choices[0].Message.Content)
	}

	return nil
}
