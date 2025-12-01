package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Client handles git operations
type Client struct{}

// NewClient creates a new git client
func NewClient() *Client {
	return &Client{}
}

// IsRepository checks if current directory is a git repository
func (c *Client) IsRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// GetDiffStaged returns the diff of staged changes
func (c *Client) GetDiffStaged() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}
	return string(output), nil
}

// GetDiffAll returns the diff of all changes (staged + unstaged)
func (c *Client) GetDiffAll() (string, error) {
	cmd := exec.Command("git", "diff", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff: %w", err)
	}
	return string(output), nil
}

// Commit creates a commit with the given message
func (c *Client) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("commit failed: %s", stderr.String())
	}
	return nil
}

// GetStatus returns the current git status
func (c *Client) GetStatus() (string, error) {
	cmd := exec.Command("git", "status", "--short")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}
	return string(output), nil
}

// HasStagedChanges checks if there are any staged changes
func (c *Client) HasStagedChanges() bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	// Returns non-zero exit code if there are staged changes
	return cmd.Run() != nil
}

// GetBranch returns the current branch name
func (c *Client) GetBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
