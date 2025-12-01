package ai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// parseSSEStream parses Server-Sent Events stream and sends chunks to the channel
func parseSSEStream(reader io.Reader, chunkChan chan<- StreamChunk) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for end of stream
		if line == "data: [DONE]" {
			chunkChan <- StreamChunk{Done: true}
			return nil
		}

		// Parse data line
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			var chunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				// Skip malformed chunks
				continue
			}

			if len(chunk.Choices) > 0 {
				choice := chunk.Choices[0]
				chunkChan <- StreamChunk{
					Content:      choice.Delta.Content,
					FinishReason: choice.FinishReason,
					Done:         choice.FinishReason != "",
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stream reading error: %w", err)
	}

	return nil
}
