package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zarazaex69/zik/apps/cli/internal/ai"
	"github.com/zarazaex69/zik/apps/cli/internal/config"
	"github.com/zarazaex69/zik/apps/cli/internal/git"
	"github.com/zarazaex69/zik/apps/cli/internal/prompt"
)

var (
	commitStaged bool
	commitAll    bool
	commitApply  bool
	commitType   string

	commitCmd = &cobra.Command{
		Use:   "commit",
		Short: "Generate conventional commit message from git diff",
		Long: `Analyze staged changes and generate a conventional commit message.
Uses AI to understand the changes and create a meaningful commit message following the Conventional Commits standard.`,
		Example: `  zik commit                    # Generate message for staged changes
  zik commit --all              # Generate message for all changes
  zik commit --apply            # Generate and apply commit
  zik commit --type feat        # Prefer 'feat' type`,
		RunE: runCommit,
	}
)

func init() {
	commitCmd.Flags().BoolVarP(&commitStaged, "staged", "s", true, "Analyze staged changes only")
	commitCmd.Flags().BoolVarP(&commitAll, "all", "a", false, "Analyze all changes (staged + unstaged)")
	commitCmd.Flags().BoolVarP(&commitApply, "apply", "y", false, "Automatically apply the generated commit message")
	commitCmd.Flags().StringVarP(&commitType, "type", "t", "", "Preferred commit type (feat, fix, docs, etc.)")
}

func runCommit(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize git client
	gitClient := git.NewClient()

	// Check if we're in a git repository
	if !gitClient.IsRepository() {
		return fmt.Errorf("not a git repository")
	}

	// Get diff based on flags
	var diff string
	if commitAll {
		diff, err = gitClient.GetDiffAll()
	} else {
		diff, err = gitClient.GetDiffStaged()
	}
	if err != nil {
		return fmt.Errorf("failed to get git diff: %w", err)
	}

	if diff == "" {
		return fmt.Errorf("no changes to commit")
	}

	// Generate commit message using AI
	fmt.Println("Analyzing changes...")

	aiClient := ai.NewClient(cfg)
	ctx := context.Background()

	// Build prompt for commit message generation
	systemPrompt := prompt.CommitSystemPrompt(cfg.Commit.ConventionalCommits, commitType)
	userPrompt := prompt.CommitUserPrompt(diff)

	for {
		messages := []ai.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		}

		// Get AI response
		resp, err := aiClient.Chat(ctx, messages, 0.3, 200) // Low temperature for consistency
		if err != nil {
			return fmt.Errorf("AI request failed: %w", err)
		}

		if len(resp.Choices) == 0 {
			return fmt.Errorf("no response from AI")
		}

		commitMessage := resp.Choices[0].Message.Content

		// Display generated commit message
		fmt.Println("\nGenerated commit message:")
		fmt.Println("─────────────────────────────")
		fmt.Println(commitMessage)
		fmt.Println("─────────────────────────────")

		// Auto-apply if flag is set
		if commitApply {
			if err := gitClient.Commit(commitMessage); err != nil {
				return fmt.Errorf("failed to commit: %w", err)
			}
			fmt.Println("\nCommit applied successfully!")
			return nil
		}

		// Interactive prompt
		fmt.Print("\nAccept? [Y]es / [N]o / [R]egenerate: ")
		var response string
		fmt.Scanln(&response)

		switch strings.ToLower(strings.TrimSpace(response)) {
		case "y", "yes", "":
			if err := gitClient.Commit(commitMessage); err != nil {
				return fmt.Errorf("failed to commit: %w", err)
			}
			fmt.Println("Commit applied successfully!")
			return nil
		case "n", "no":
			fmt.Println("Commit cancelled.")
			return nil
		case "r", "regenerate":
			fmt.Println("\nRegenerating...")
			continue
		default:
			fmt.Println("Invalid option. Please choose Y, N, or R.")
			continue
		}
	}
}
