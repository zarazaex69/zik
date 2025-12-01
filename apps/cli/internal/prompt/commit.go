package prompt

import "fmt"

// CommitSystemPrompt generates the system prompt for commit message generation
func CommitSystemPrompt(conventional bool, preferredType string) string {
	base := `You are an expert at analyzing code changes and writing clear, concise commit messages.
Your task is to analyze git diff output and generate a meaningful commit message.`

	if conventional {
		base += `

Follow the Conventional Commits standard strictly:
- Format: type(scope): description
- Types: feat, fix, docs, style, refactor, perf, test, chore
- Scope: optional, indicates the area of change
- Description: imperative mood, lowercase, no period at end

Examples:
- feat(auth): add JWT token validation
- fix(api): handle null response in user endpoint
- docs(readme): update installation instructions
- refactor(parser): simplify token extraction logic`

		if preferredType != "" {
			base += fmt.Sprintf("\n- Prefer using '%s' type when appropriate", preferredType)
		}
	}

	base += `

Rules:
1. Keep the message concise (max 72 characters for first line)
2. Focus on WHAT changed and WHY, not HOW
3. Use imperative mood ("add" not "added" or "adds")
4. Be specific but brief
5. Return ONLY the commit message, no explanations or additional text`

	return base
}

// CommitUserPrompt generates the user prompt with the git diff
func CommitUserPrompt(diff string) string {
	return fmt.Sprintf(`Analyze the following git diff and generate a commit message:

%s

Generate a commit message for these changes.`, diff)
}
