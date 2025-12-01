package render

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles for different markdown elements
	boldStyle      = lipgloss.NewStyle().Bold(true)
	italicStyle    = lipgloss.NewStyle().Italic(true)
	codeStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	h1Style        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	h2Style        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	h3Style        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226"))
	linkStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("33")).Underline(true)
	linkURLStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	codeBlockStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	dimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// MarkdownRenderer handles markdown formatting for terminal output with streaming support
type MarkdownRenderer struct {
	inCodeBlock   bool
	codeBlockLang string
	lineBuffer    string
}

// NewMarkdownRenderer creates a new markdown renderer
func NewMarkdownRenderer() *MarkdownRenderer {
	return &MarkdownRenderer{}
}

// ProcessChunk processes a chunk of markdown text and returns formatted output
// This works incrementally for streaming - it buffers incomplete lines
func (r *MarkdownRenderer) ProcessChunk(chunk string) string {
	if chunk == "" {
		return ""
	}

	var result strings.Builder
	r.lineBuffer += chunk

	// Process complete lines
	lines := strings.Split(r.lineBuffer, "\n")

	// Keep the last incomplete line in buffer
	if !strings.HasSuffix(r.lineBuffer, "\n") {
		r.lineBuffer = lines[len(lines)-1]
		lines = lines[:len(lines)-1]
	} else {
		r.lineBuffer = ""
	}

	// Format each complete line
	for i, line := range lines {
		formatted := r.formatLine(line)
		result.WriteString(formatted)
		if i < len(lines)-1 || len(lines) > 0 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// Flush returns any remaining buffered content
func (r *MarkdownRenderer) Flush() string {
	if r.lineBuffer == "" {
		return ""
	}
	formatted := r.formatLine(r.lineBuffer)
	r.lineBuffer = ""
	return formatted
}

// formatLine formats a single line of markdown
func (r *MarkdownRenderer) formatLine(line string) string {
	// Check for code block markers
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "```") {
		if !r.inCodeBlock {
			r.inCodeBlock = true
			lang := strings.TrimPrefix(trimmed, "```")
			r.codeBlockLang = lang
			if lang != "" {
				return dimStyle.Render("╭─ " + lang)
			}
			return dimStyle.Render("╭─ code")
		} else {
			r.inCodeBlock = false
			r.codeBlockLang = ""
			return dimStyle.Render("╰─")
		}
	}

	// If inside code block, just add prefix and style
	if r.inCodeBlock {
		return dimStyle.Render("│ ") + codeBlockStyle.Render(line)
	}

	// Format headers
	if strings.HasPrefix(line, "### ") {
		text := strings.TrimPrefix(line, "### ")
		return h3Style.Render(text)
	}
	if strings.HasPrefix(line, "## ") {
		text := strings.TrimPrefix(line, "## ")
		return h2Style.Render(text)
	}
	if strings.HasPrefix(line, "# ") {
		text := strings.TrimPrefix(line, "# ")
		return h1Style.Render(text)
	}

	// Format inline elements
	line = r.formatInline(line)

	return line
}

// formatInline formats inline markdown elements
func (r *MarkdownRenderer) formatInline(text string) string {
	// Process in order: code, bold, italic, links
	// This prevents nested formatting issues

	// Inline code: `code` - do this first to protect code content
	codeRe := regexp.MustCompile("`([^`]+)`")
	codeParts := codeRe.FindAllStringSubmatchIndex(text, -1)

	var result strings.Builder
	lastIdx := 0

	for _, match := range codeParts {
		// Add text before code
		beforeCode := text[lastIdx:match[0]]
		result.WriteString(r.formatTextWithoutCode(beforeCode))

		// Add formatted code
		codeContent := text[match[2]:match[3]]
		result.WriteString(codeStyle.Render(codeContent))

		lastIdx = match[1]
	}

	// Add remaining text
	if lastIdx < len(text) {
		result.WriteString(r.formatTextWithoutCode(text[lastIdx:]))
	}

	if result.Len() > 0 {
		return result.String()
	}

	return r.formatTextWithoutCode(text)
}

// formatTextWithoutCode formats text that doesn't contain inline code
func (r *MarkdownRenderer) formatTextWithoutCode(text string) string {
	// Bold: **text** or __text__
	boldRe := regexp.MustCompile(`\*\*([^*]+)\*\*|__([^_]+)__`)
	text = boldRe.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(match, "*_")
		return boldStyle.Render(content)
	})

	// Italic: *text* or _text_
	italicRe := regexp.MustCompile(`\*([^*\s][^*]*?)\*|_([^_\s][^_]*?)_`)
	text = italicRe.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(match, "*_")
		return italicStyle.Render(content)
	})

	// Links: [text](url)
	linkRe := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = linkRe.ReplaceAllStringFunc(text, func(match string) string {
		parts := linkRe.FindStringSubmatch(match)
		if len(parts) == 3 {
			return linkStyle.Render(parts[1]) + " " + linkURLStyle.Render("("+parts[2]+")")
		}
		return match
	})

	return text
}

// Reset resets the renderer state
func (r *MarkdownRenderer) Reset() {
	r.inCodeBlock = false
	r.codeBlockLang = ""
	r.lineBuffer = ""
}
