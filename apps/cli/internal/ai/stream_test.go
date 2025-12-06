package ai

import (
	"strings"
	"testing"
)

func TestParseSSEStream_EmptyStream(t *testing.T) {
	reader := strings.NewReader("")
	chunkChan := make(chan StreamChunk, 10)

	err := parseSSEStream(reader, chunkChan)
	if err != nil {
		t.Errorf("parseSSEStream() error = %v, want nil for empty stream", err)
	}

	close(chunkChan)

	// Should not receive any chunks
	count := 0
	for range chunkChan {
		count++
	}

	if count != 0 {
		t.Errorf("parseSSEStream() sent %d chunks, want 0 for empty stream", count)
	}
}

func TestParseSSEStream_DoneMessage(t *testing.T) {
	input := "data: [DONE]\n"
	reader := strings.NewReader(input)
	chunkChan := make(chan StreamChunk, 10)

	err := parseSSEStream(reader, chunkChan)
	if err != nil {
		t.Errorf("parseSSEStream() error = %v", err)
	}

	close(chunkChan)

	// Should receive one Done chunk
	chunks := make([]StreamChunk, 0)
	for chunk := range chunkChan {
		chunks = append(chunks, chunk)
	}

	if len(chunks) != 1 {
		t.Fatalf("parseSSEStream() sent %d chunks, want 1", len(chunks))
	}

	if !chunks[0].Done {
		t.Error("parseSSEStream() chunk.Done = false, want true")
	}
}

func TestParseSSEStream_ValidChunk(t *testing.T) {
	input := `data: {"choices":[{"delta":{"content":"Hello"},"finish_reason":""}]}

`
	reader := strings.NewReader(input)
	chunkChan := make(chan StreamChunk, 10)

	err := parseSSEStream(reader, chunkChan)
	if err != nil {
		t.Errorf("parseSSEStream() error = %v", err)
	}

	close(chunkChan)

	chunks := make([]StreamChunk, 0)
	for chunk := range chunkChan {
		chunks = append(chunks, chunk)
	}

	if len(chunks) != 1 {
		t.Fatalf("parseSSEStream() sent %d chunks, want 1", len(chunks))
	}

	if chunks[0].Content != "Hello" {
		t.Errorf("parseSSEStream() chunk.Content = %q, want Hello", chunks[0].Content)
	}

	if chunks[0].Done {
		t.Error("parseSSEStream() chunk.Done = true, want false")
	}
}

func TestParseSSEStream_FinishReason(t *testing.T) {
	input := `data: {"choices":[{"delta":{"content":""},"finish_reason":"stop"}]}

`
	reader := strings.NewReader(input)
	chunkChan := make(chan StreamChunk, 10)

	err := parseSSEStream(reader, chunkChan)
	if err != nil {
		t.Errorf("parseSSEStream() error = %v", err)
	}

	close(chunkChan)

	chunks := make([]StreamChunk, 0)
	for chunk := range chunkChan {
		chunks = append(chunks, chunk)
	}

	if len(chunks) != 1 {
		t.Fatalf("parseSSEStream() sent %d chunks, want 1", len(chunks))
	}

	if chunks[0].FinishReason != "stop" {
		t.Errorf("parseSSEStream() chunk.FinishReason = %q, want stop", chunks[0].FinishReason)
	}

	if !chunks[0].Done {
		t.Error("parseSSEStream() chunk.Done = false, want true when finish_reason is set")
	}
}

func TestParseSSEStream_MultipleChunks(t *testing.T) {
	input := `data: {"choices":[{"delta":{"content":"Hello"},"finish_reason":""}]}

data: {"choices":[{"delta":{"content":" world"},"finish_reason":""}]}

data: {"choices":[{"delta":{"content":"!"},"finish_reason":"stop"}]}

`
	reader := strings.NewReader(input)
	chunkChan := make(chan StreamChunk, 10)

	err := parseSSEStream(reader, chunkChan)
	if err != nil {
		t.Errorf("parseSSEStream() error = %v", err)
	}

	close(chunkChan)

	chunks := make([]StreamChunk, 0)
	for chunk := range chunkChan {
		chunks = append(chunks, chunk)
	}

	if len(chunks) != 3 {
		t.Fatalf("parseSSEStream() sent %d chunks, want 3", len(chunks))
	}

	expectedContent := []string{"Hello", " world", "!"}
	for i, expected := range expectedContent {
		if chunks[i].Content != expected {
			t.Errorf("chunk[%d].Content = %q, want %q", i, chunks[i].Content, expected)
		}
	}

	// Only last chunk should be done
	if chunks[0].Done || chunks[1].Done {
		t.Error("Early chunks should not be done")
	}
	if !chunks[2].Done {
		t.Error("Last chunk should be done")
	}
}

func TestParseSSEStream_MalformedJSON(t *testing.T) {
	input := `data: {invalid json}

data: {"choices":[{"delta":{"content":"valid"},"finish_reason":""}]}

`
	reader := strings.NewReader(input)
	chunkChan := make(chan StreamChunk, 10)

	err := parseSSEStream(reader, chunkChan)
	if err != nil {
		t.Errorf("parseSSEStream() error = %v, should skip malformed chunks", err)
	}

	close(chunkChan)

	chunks := make([]StreamChunk, 0)
	for chunk := range chunkChan {
		chunks = append(chunks, chunk)
	}

	// Should only receive the valid chunk
	if len(chunks) != 1 {
		t.Fatalf("parseSSEStream() sent %d chunks, want 1 (malformed should be skipped)", len(chunks))
	}

	if chunks[0].Content != "valid" {
		t.Errorf("parseSSEStream() chunk.Content = %q, want valid", chunks[0].Content)
	}
}

func TestParseSSEStream_EmptyLines(t *testing.T) {
	input := `

data: {"choices":[{"delta":{"content":"test"},"finish_reason":""}]}


data: [DONE]
`
	reader := strings.NewReader(input)
	chunkChan := make(chan StreamChunk, 10)

	err := parseSSEStream(reader, chunkChan)
	if err != nil {
		t.Errorf("parseSSEStream() error = %v", err)
	}

	close(chunkChan)

	chunks := make([]StreamChunk, 0)
	for chunk := range chunkChan {
		chunks = append(chunks, chunk)
	}

	// Should receive content chunk and done chunk
	if len(chunks) != 2 {
		t.Fatalf("parseSSEStream() sent %d chunks, want 2", len(chunks))
	}
}

func TestParseSSEStream_NoChoices(t *testing.T) {
	input := `data: {"choices":[]}

`
	reader := strings.NewReader(input)
	chunkChan := make(chan StreamChunk, 10)

	err := parseSSEStream(reader, chunkChan)
	if err != nil {
		t.Errorf("parseSSEStream() error = %v", err)
	}

	close(chunkChan)

	chunks := make([]StreamChunk, 0)
	for chunk := range chunkChan {
		chunks = append(chunks, chunk)
	}

	// Should not send chunk if choices array is empty
	if len(chunks) != 0 {
		t.Errorf("parseSSEStream() sent %d chunks, want 0 for empty choices", len(chunks))
	}
}

func TestParseSSEStream_NonDataLines(t *testing.T) {
	input := `: comment line
event: message
data: {"choices":[{"delta":{"content":"test"},"finish_reason":""}]}

`
	reader := strings.NewReader(input)
	chunkChan := make(chan StreamChunk, 10)

	err := parseSSEStream(reader, chunkChan)
	if err != nil {
		t.Errorf("parseSSEStream() error = %v", err)
	}

	close(chunkChan)

	chunks := make([]StreamChunk, 0)
	for chunk := range chunkChan {
		chunks = append(chunks, chunk)
	}

	// Should only process data: lines
	if len(chunks) != 1 {
		t.Fatalf("parseSSEStream() sent %d chunks, want 1", len(chunks))
	}

	if chunks[0].Content != "test" {
		t.Errorf("parseSSEStream() chunk.Content = %q, want test", chunks[0].Content)
	}
}
