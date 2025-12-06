package git

import (
	"os"
	"os/exec"
	"testing"
)

// setupTestRepo creates a temporary git repository for testing
func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to init git repo: %v", err)
	}

	// Configure git user for commits
	exec.Command("git", "config", "user.email", "test@example.com").Dir = tmpDir
	exec.Command("git", "config", "user.name", "Test User").Dir = tmpDir

	// Change to test directory
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)

	cleanup := func() {
		os.Chdir(oldDir)
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

func TestIsRepository(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) func()
		expected bool
	}{
		{
			name: "valid git repository",
			setup: func(t *testing.T) func() {
				_, cleanup := setupTestRepo(t)
				return cleanup
			},
			expected: true,
		},
		{
			name: "not a git repository",
			setup: func(t *testing.T) func() {
				tmpDir, err := os.MkdirTemp("", "not-git-*")
				if err != nil {
					t.Fatalf("failed to create temp dir: %v", err)
				}
				oldDir, _ := os.Getwd()
				os.Chdir(tmpDir)
				return func() {
					os.Chdir(oldDir)
					os.RemoveAll(tmpDir)
				}
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup(t)
			defer cleanup()

			client := NewClient()
			result := client.IsRepository()

			if result != tt.expected {
				t.Errorf("IsRepository() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetDiffStaged(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	client := NewClient()

	// Create and stage a file
	testFile := "test.txt"
	content := []byte("test content")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to stage file: %v", err)
	}

	diff, err := client.GetDiffStaged()
	if err != nil {
		t.Fatalf("GetDiffStaged() error = %v", err)
	}

	// Diff should contain the filename
	if diff == "" {
		t.Error("GetDiffStaged() returned empty diff for staged changes")
	}
}

func TestGetDiffAll(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	client := NewClient()

	// Create initial commit
	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("initial"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	exec.Command("git", "add", testFile).Run()
	exec.Command("git", "commit", "-m", "initial").Run()

	// Modify file
	if err := os.WriteFile(testFile, []byte("modified"), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	diff, err := client.GetDiffAll()
	if err != nil {
		t.Fatalf("GetDiffAll() error = %v", err)
	}

	// Diff should show changes
	if diff == "" {
		t.Error("GetDiffAll() returned empty diff for modified file")
	}
}

func TestCommit(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	client := NewClient()

	// Create and stage a file
	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	exec.Command("git", "add", testFile).Run()

	// Test commit
	err := client.Commit("test: add test file")
	if err != nil {
		t.Errorf("Commit() error = %v", err)
	}

	// Verify commit was created
	cmd := exec.Command("git", "log", "--oneline")
	output, _ := cmd.Output()
	if len(output) == 0 {
		t.Error("Commit() did not create a commit")
	}
}

func TestCommit_NoStagedChanges(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	client := NewClient()

	// Try to commit without staged changes
	err := client.Commit("test: empty commit")
	if err == nil {
		t.Error("Commit() should fail with no staged changes")
	}
}

func TestGetStatus(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	client := NewClient()

	// Create untracked file
	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	status, err := client.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}

	// Status should show untracked file
	if status == "" {
		t.Error("GetStatus() returned empty status for untracked file")
	}
}

func TestHasStagedChanges(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	client := NewClient()

	// Initially no staged changes
	if client.HasStagedChanges() {
		t.Error("HasStagedChanges() = true, want false for empty repo")
	}

	// Create and stage a file
	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	exec.Command("git", "add", testFile).Run()

	// Now should have staged changes
	if !client.HasStagedChanges() {
		t.Error("HasStagedChanges() = false, want true after staging file")
	}
}

func TestGetBranch(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	client := NewClient()

	// Create initial commit to establish branch
	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	exec.Command("git", "add", testFile).Run()
	exec.Command("git", "commit", "-m", "initial").Run()

	branch, err := client.GetBranch()
	if err != nil {
		t.Fatalf("GetBranch() error = %v", err)
	}

	// Default branch should be master or main
	if branch != "master" && branch != "main" {
		t.Errorf("GetBranch() = %v, want master or main", branch)
	}
}
