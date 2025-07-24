package provider

import (
	"strings"
	"testing"
)

func TestToolCallMarkerFiltering(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool // true if content should be filtered out
	}{
		{
			name:     "normal content passes through",
			content:  "Hello, how can I help you?",
			expected: false,
		},
		{
			name:     "content with tool begin marker is filtered",
			content:  "Some text <|tool_calls_section_begin|> more text",
			expected: true,
		},
		{
			name:     "content with tool end marker is filtered",
			content:  "Some text <|tool_calls_section_end|> more text",
			expected: true,
		},
		{
			name:     "content with both markers is filtered",
			content:  "<|tool_calls_section_begin|>tool call<|tool_calls_section_end|>",
			expected: true,
		},
		{
			name:     "empty content is not filtered",
			content:  "",
			expected: false,
		},
		{
			name:     "similar but different markers pass through",
			content:  "<|other_section_begin|>some content<|other_section_end|>",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the filtering logic
			shouldFilter := strings.Contains(tt.content, ToolBegin) || strings.Contains(tt.content, ToolEnd)
			
			if shouldFilter != tt.expected {
				t.Errorf("Expected filtering decision %v for content %q, got %v", 
					tt.expected, tt.content, shouldFilter)
			}
		})
	}
}

func TestToolCallMarkerConstants(t *testing.T) {
	// Verify the constants are defined correctly
	if ToolBegin != "<|tool_calls_section_begin|>" {
		t.Errorf("ToolBegin constant is incorrect: got %q", ToolBegin)
	}
	
	if ToolEnd != "<|tool_calls_section_end|>" {
		t.Errorf("ToolEnd constant is incorrect: got %q", ToolEnd)
	}
}

func TestProviderEventFiltering(t *testing.T) {
	// Test that demonstrates how the filtering would work in the streaming context
	testCases := []struct {
		name                 string
		deltaContent        string
		shouldGenerateEvent bool
		expectedContent     string
	}{
		{
			name:                 "normal content generates event",
			deltaContent:        "Hello world",
			shouldGenerateEvent: true,
			expectedContent:     "Hello world",
		},
		{
			name:                 "content with tool begin marker does not generate event",
			deltaContent:        "text <|tool_calls_section_begin|> more",
			shouldGenerateEvent: false,
			expectedContent:     "",
		},
		{
			name:                 "content with tool end marker does not generate event", 
			deltaContent:        "<|tool_calls_section_end|> after tool",
			shouldGenerateEvent: false,
			expectedContent:     "",
		},
		{
			name:                 "empty content does not generate event",
			deltaContent:        "",
			shouldGenerateEvent: false,
			expectedContent:     "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the filtering logic from both openai.go and copilot.go
			var event *ProviderEvent
			var accumulatedContent string

			if tc.deltaContent != "" {
				content := tc.deltaContent
				if !strings.Contains(content, ToolBegin) && !strings.Contains(content, ToolEnd) {
					event = &ProviderEvent{
						Type:    EventContentDelta,
						Content: content,
					}
					accumulatedContent += content
				}
			}

			if tc.shouldGenerateEvent {
				if event == nil {
					t.Errorf("Expected event to be generated for content %q, but got nil", tc.deltaContent)
				} else if event.Content != tc.expectedContent {
					t.Errorf("Expected event content %q, got %q", tc.expectedContent, event.Content)
				}
			} else {
				if event != nil {
					t.Errorf("Expected no event for content %q, but got event with content %q", 
						tc.deltaContent, event.Content)
				}
			}

			if accumulatedContent != tc.expectedContent {
				t.Errorf("Expected accumulated content %q, got %q", tc.expectedContent, accumulatedContent)
			}
		})
	}
}
