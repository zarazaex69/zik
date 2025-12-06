package render

import (
	"strings"
	"testing"
)

func TestNewMarkdownRenderer(t *testing.T) {
	renderer := NewMarkdownRenderer()
	if renderer == nil {
		t.Fatal("NewMarkdownRenderer() returned nil")
	}

	if renderer.inCodeBlock {
		t.Error("New renderer should not be in code block")
	}
	if renderer.lineBuffer != "" {
		t.Error("New renderer should have empty line buffer")
	}
}

func TestProcessChunk_EmptyInput(t *testing.T) {
	renderer := NewMarkdownRenderer()
	result := renderer.ProcessChunk("")

	if result != "" {
		t.Errorf("ProcessChunk(\"\") = %q, want empty string", result)
	}
}

func TestProcessChunk_SimpleText(t *testing.T) {
	renderer := NewMarkdownRenderer()
	input := "Hello, world!\n"
	result := renderer.ProcessChunk(input)

	if result == "" {
		t.Error("ProcessChunk() returned empty string for simple text")
	}
}

func TestProcessChunk_Headers(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"H1", "# Header 1\n"},
		{"H2", "## Header 2\n"},
		{"H3", "### Header 3\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewMarkdownRenderer()
			result := renderer.ProcessChunk(tt.input)

			if result == "" {
				t.Errorf("ProcessChunk(%q) returned empty string", tt.input)
			}
		})
	}
}

func TestProcessChunk_CodeBlock(t *testing.T) {
	renderer := NewMarkdownRenderer()

	// Start code block
	result1 := renderer.ProcessChunk("```go\n")
	if result1 == "" {
		t.Error("ProcessChunk() should return formatted code block start")
	}

	if !renderer.inCodeBlock {
		t.Error("Renderer should be in code block after opening marker")
	}

	// Code content
	result2 := renderer.ProcessChunk("func main() {\n")
	if result2 == "" {
		t.Error("ProcessChunk() should return formatted code line")
	}

	// End code block
	result3 := renderer.ProcessChunk("```\n")
	if result3 == "" {
		t.Error("ProcessChunk() should return formatted code block end")
	}

	if renderer.inCodeBlock {
		t.Error("Renderer should not be in code block after closing marker")
	}
}

func TestProcessChunk_InlineCode(t *testing.T) {
	renderer := NewMarkdownRenderer()
	input := "Use `fmt.Println()` to print\n"
	result := renderer.ProcessChunk(input)

	if result == "" {
		t.Error("ProcessChunk() returned empty string for inline code")
	}
}

func TestProcessChunk_Bold(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Asterisks", "This is **bold** text\n"},
		{"Underscores", "This is __bold__ text\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewMarkdownRenderer()
			result := renderer.ProcessChunk(tt.input)

			if result == "" {
				t.Errorf("ProcessChunk(%q) returned empty string", tt.input)
			}
		})
	}
}

func TestProcessChunk_Italic(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Asterisk", "This is *italic* text\n"},
		{"Underscore", "This is _italic_ text\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewMarkdownRenderer()
			result := renderer.ProcessChunk(tt.input)

			if result == "" {
				t.Errorf("ProcessChunk(%q) returned empty string", tt.input)
			}
		})
	}
}

func TestProcessChunk_Links(t *testing.T) {
	renderer := NewMarkdownRenderer()
	input := "Check [documentation](https://example.com)\n"
	result := renderer.ProcessChunk(input)

	if result == "" {
		t.Error("ProcessChunk() returned empty string for link")
	}
}

func TestProcessChunk_Streaming(t *testing.T) {
	renderer := NewMarkdownRenderer()

	// Simulate streaming: incomplete line
	result1 := renderer.ProcessChunk("Hello")
	if result1 != "" {
		t.Error("ProcessChunk() should buffer incomplete line")
	}

	// Complete the line
	result2 := renderer.ProcessChunk(" world\n")
	if result2 == "" {
		t.Error("ProcessChunk() should return complete line")
	}
}

func TestFlush(t *testing.T) {
	renderer := NewMarkdownRenderer()

	// Add incomplete line to buffer
	renderer.ProcessChunk("Incomplete line")

	// Flush should return buffered content
	result := renderer.Flush()
	if result == "" {
		t.Error("Flush() should return buffered content")
	}

	// Buffer should be empty after flush
	if renderer.lineBuffer != "" {
		t.Error("Flush() should clear line buffer")
	}

	// Second flush should return empty
	result2 := renderer.Flush()
	if result2 != "" {
		t.Error("Flush() on empty buffer should return empty string")
	}
}

func TestReset(t *testing.T) {
	renderer := NewMarkdownRenderer()

	// Set some state
	renderer.ProcessChunk("```go\n")
	renderer.ProcessChunk("code")

	// Reset
	renderer.Reset()

	if renderer.inCodeBlock {
		t.Error("Reset() should clear inCodeBlock flag")
	}
	if renderer.codeBlockLang != "" {
		t.Error("Reset() should clear codeBlockLang")
	}
	if renderer.lineBuffer != "" {
		t.Error("Reset() should clear lineBuffer")
	}
}

func TestFormatLine_CodeBlockToggle(t *testing.T) {
	renderer := NewMarkdownRenderer()

	// Test opening code block
	result := renderer.formatLine("```python")
	if !renderer.inCodeBlock {
		t.Error("formatLine() should set inCodeBlock to true")
	}
	if renderer.codeBlockLang != "python" {
		t.Errorf("formatLine() codeBlockLang = %q, want python", renderer.codeBlockLang)
	}
	if result == "" {
		t.Error("formatLine() should return formatted opening marker")
	}

	// Test closing code block
	result = renderer.formatLine("```")
	if renderer.inCodeBlock {
		t.Error("formatLine() should set inCodeBlock to false")
	}
	if result == "" {
		t.Error("formatLine() should return formatted closing marker")
	}
}

func TestFormatInline_MultipleElements(t *testing.T) {
	renderer := NewMarkdownRenderer()
	input := "Use **bold**, *italic*, and `code` together"
	result := renderer.formatInline(input)

	if result == "" {
		t.Error("formatInline() returned empty string for mixed formatting")
	}
}

func TestFormatTextWithoutCode_BoldAndItalic(t *testing.T) {
	renderer := NewMarkdownRenderer()

	tests := []struct {
		name  string
		input string
	}{
		{"Bold with asterisks", "**bold**"},
		{"Bold with underscores", "__bold__"},
		{"Italic with asterisk", "*italic*"},
		{"Italic with underscore", "_italic_"},
		{"Mixed", "**bold** and *italic*"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderer.formatTextWithoutCode(tt.input)
			if result == "" {
				t.Errorf("formatTextWithoutCode(%q) returned empty string", tt.input)
			}
		})
	}
}

func TestProcessChunk_ComplexMarkdown(t *testing.T) {
	renderer := NewMarkdownRenderer()

	markdown := `# Title

This is a paragraph with **bold** and *italic* text.

## Code Example

` + "```go" + `
func main() {
    fmt.Println("Hello")
}
` + "```" + `

Check [this link](https://example.com) for more info.
`

	lines := strings.Split(markdown, "\n")
	for _, line := range lines {
		renderer.ProcessChunk(line + "\n")
	}

	// Should complete without errors
	final := renderer.Flush()
	_ = final // Just verify no panic
}

func TestProcessChunk_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty line", "\n"},
		{"Only spaces", "   \n"},
		{"Multiple newlines", "\n\n\n"},
		{"Tab characters", "\t\tcode\n"},
		{"Mixed whitespace", "  \t  text  \t  \n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewMarkdownRenderer()
			// Should not panic
			_ = renderer.ProcessChunk(tt.input)
		})
	}
}
