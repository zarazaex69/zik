package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zarazaex69/zik/apps/cli/internal/ai"
	"github.com/zarazaex69/zik/apps/cli/internal/config"
	"github.com/zarazaex69/zik/apps/cli/internal/prompt"
	"github.com/zarazaex69/zik/apps/cli/internal/render"
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

	// Build messages with system prompt to constrain formatting
	messages := []ai.Message{
		{Role: "system", Content: prompt.AskSystemPrompt()},
		{Role: "user", Content: question},
	}

	if askStream {
		// Streaming response with markdown rendering
		renderer := render.NewMarkdownRenderer()
		chunkChan, errChan := aiClient.ChatStream(ctx, messages, cfg.Temperature, cfg.MaxTokens)

		for {
			select {
			case chunk, ok := <-chunkChan:
				if !ok {
					// Flush any remaining buffered content
					if remaining := renderer.Flush(); remaining != "" {
						fmt.Print(remaining)
					}
					fmt.Println() // New line at end
					return nil
				}
				if chunk.Content != "" {
					formatted := renderer.ProcessChunk(chunk.Content)
					fmt.Print(formatted)
				}
			case err := <-errChan:
				if err != nil {
					return fmt.Errorf("AI request failed: %w", err)
				}
			}
		}
	} else {
		// Non-streaming response with markdown rendering
		resp, err := aiClient.Chat(ctx, messages, cfg.Temperature, cfg.MaxTokens)
		if err != nil {
			return fmt.Errorf("AI request failed: %w", err)
		}

		if len(resp.Choices) == 0 {
			return fmt.Errorf("no response from AI")
		}

		renderer := render.NewMarkdownRenderer()
		formatted := renderer.ProcessChunk(resp.Choices[0].Message.Content)
		fmt.Println(formatted)
	}

	return nil
}
