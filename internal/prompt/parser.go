package prompt

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const frontmatterDelimiter = "---"

// Parse parses markdown content with YAML frontmatter
func Parse(content string) (*Prompt, error) {
	frontmatter, body, err := splitFrontmatter(content)
	if err != nil {
		return nil, err
	}

	var p Prompt
	if frontmatter != "" {
		if err := yaml.Unmarshal([]byte(frontmatter), &p); err != nil {
			return nil, fmt.Errorf("invalid YAML frontmatter: %w", err)
		}
	}

	p.Content = strings.TrimSpace(body)

	return &p, nil
}

// splitFrontmatter separates the YAML frontmatter from the content
func splitFrontmatter(content string) (frontmatter, body string, err error) {
	content = strings.TrimSpace(content)

	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(content, frontmatterDelimiter) {
		return "", content, nil
	}

	// Find the closing delimiter
	rest := content[len(frontmatterDelimiter):]

	// Handle newline after opening delimiter
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	} else if len(rest) > 1 && rest[0] == '\r' && rest[1] == '\n' {
		rest = rest[2:]
	}

	// Find closing delimiter
	endIdx := strings.Index(rest, "\n"+frontmatterDelimiter)
	if endIdx == -1 {
		// Try with \r\n
		endIdx = strings.Index(rest, "\r\n"+frontmatterDelimiter)
		if endIdx == -1 {
			return "", content, fmt.Errorf("missing closing frontmatter delimiter")
		}
	}

	frontmatter = rest[:endIdx]
	body = rest[endIdx+1+len(frontmatterDelimiter):]

	// Trim leading newline from body
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	} else if len(body) > 1 && body[0] == '\r' && body[1] == '\n' {
		body = body[2:]
	}

	return frontmatter, body, nil
}

// HasFrontmatter checks if content has valid frontmatter
func HasFrontmatter(content string) bool {
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, frontmatterDelimiter) {
		return false
	}

	rest := content[len(frontmatterDelimiter):]
	return strings.Contains(rest, "\n"+frontmatterDelimiter)
}

// ParseFrontmatterOnly parses just the frontmatter without the content
func ParseFrontmatterOnly(content string) (*Prompt, error) {
	frontmatter, _, err := splitFrontmatter(content)
	if err != nil {
		return nil, err
	}

	if frontmatter == "" {
		return nil, fmt.Errorf("no frontmatter found")
	}

	var p Prompt
	if err := yaml.Unmarshal([]byte(frontmatter), &p); err != nil {
		return nil, fmt.Errorf("invalid YAML frontmatter: %w", err)
	}

	return &p, nil
}
