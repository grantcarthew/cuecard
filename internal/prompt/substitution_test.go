package prompt

import (
	"strings"
	"testing"
	"time"
)

func TestSubstitute(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		resolver *VariableResolver
		check    func(result string) bool
		desc     string
	}{
		{
			name:    "substitute INPUT",
			content: "Hello ${INPUT}!",
			resolver: &VariableResolver{
				Input: "World",
			},
			check: func(result string) bool {
				return result == "Hello World!"
			},
			desc: "expected 'Hello World!'",
		},
		{
			name:    "substitute CLIPBOARD",
			content: "Pasted: ${CLIPBOARD}",
			resolver: &VariableResolver{
				Clipboard: "clipboard content",
			},
			check: func(result string) bool {
				return result == "Pasted: clipboard content"
			},
			desc: "expected clipboard content substituted",
		},
		{
			name:    "substitute DATE",
			content: "Today is ${DATE}",
			resolver: &VariableResolver{},
			check: func(result string) bool {
				today := time.Now().Format("2006-01-02")
				return strings.Contains(result, today)
			},
			desc: "expected today's date in ISO format",
		},
		{
			name:    "substitute DATETIME",
			content: "Now is ${DATETIME}",
			resolver: &VariableResolver{},
			check: func(result string) bool {
				// Just check it contains the date portion
				today := time.Now().Format("2006-01-02")
				return strings.Contains(result, today)
			},
			desc: "expected current datetime",
		},
		{
			name:    "multiple substitutions",
			content: "Input: ${INPUT}, Date: ${DATE}",
			resolver: &VariableResolver{
				Input: "test",
			},
			check: func(result string) bool {
				today := time.Now().Format("2006-01-02")
				return strings.Contains(result, "Input: test") && strings.Contains(result, today)
			},
			desc: "expected both INPUT and DATE substituted",
		},
		{
			name:     "nil resolver",
			content:  "Hello ${INPUT}!",
			resolver: nil,
			check: func(result string) bool {
				return result == "Hello !"
			},
			desc: "expected empty INPUT with nil resolver",
		},
		{
			name:    "no variables",
			content: "Plain text without variables",
			resolver: &VariableResolver{
				Input: "unused",
			},
			check: func(result string) bool {
				return result == "Plain text without variables"
			},
			desc: "expected unchanged content",
		},
		{
			name:    "FILE variable without selector",
			content: "File: ${FILE}",
			resolver: &VariableResolver{
				FileSelector: nil,
			},
			check: func(result string) bool {
				return result == "File: ${FILE}"
			},
			desc: "expected FILE unchanged without selector",
		},
		{
			name:    "FILE variable with selector",
			content: "File: ${FILE}",
			resolver: &VariableResolver{
				FileSelector: func() string {
					return "/path/to/file.txt"
				},
			},
			check: func(result string) bool {
				return result == "File: /path/to/file.txt"
			},
			desc: "expected file path substituted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Substitute(tt.content, tt.resolver)
			if !tt.check(result) {
				t.Errorf("Substitute() = %q, %s", result, tt.desc)
			}
		})
	}
}

func TestSubstituteWithValues(t *testing.T) {
	tests := []struct {
		name    string
		content string
		values  map[string]string
		want    string
	}{
		{
			name:    "single value",
			content: "Hello ${NAME}!",
			values:  map[string]string{"NAME": "World"},
			want:    "Hello World!",
		},
		{
			name:    "multiple values",
			content: "${GREETING} ${NAME}!",
			values:  map[string]string{"GREETING": "Hello", "NAME": "World"},
			want:    "Hello World!",
		},
		{
			name:    "unused values",
			content: "Hello!",
			values:  map[string]string{"NAME": "World"},
			want:    "Hello!",
		},
		{
			name:    "missing value",
			content: "Hello ${NAME}!",
			values:  map[string]string{},
			want:    "Hello ${NAME}!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubstituteWithValues(tt.content, tt.values); got != tt.want {
				t.Errorf("SubstituteWithValues() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractVariables(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "single variable",
			content: "Hello ${NAME}!",
			want:    []string{"NAME"},
		},
		{
			name:    "multiple variables",
			content: "${GREETING} ${NAME}!",
			want:    []string{"GREETING", "NAME"},
		},
		{
			name:    "duplicate variables",
			content: "${NAME} and ${NAME}",
			want:    []string{"NAME"},
		},
		{
			name:    "no variables",
			content: "Plain text",
			want:    nil,
		},
		{
			name:    "all standard variables",
			content: "${INPUT} ${DATE} ${DATETIME} ${CLIPBOARD} ${FILE}",
			want:    []string{"INPUT", "DATE", "DATETIME", "CLIPBOARD", "FILE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractVariables(tt.content)
			if len(got) != len(tt.want) {
				t.Errorf("ExtractVariables() = %v, want %v", got, tt.want)
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("ExtractVariables()[%d] = %q, want %q", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestContainsVariable(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		variable string
		want     bool
	}{
		{
			name:     "contains variable",
			content:  "Hello ${NAME}!",
			variable: "NAME",
			want:     true,
		},
		{
			name:     "does not contain variable",
			content:  "Hello ${NAME}!",
			variable: "OTHER",
			want:     false,
		},
		{
			name:     "empty content",
			content:  "",
			variable: "NAME",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsVariable(tt.content, tt.variable); got != tt.want {
				t.Errorf("ContainsVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsAnyVariable(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "contains variable",
			content: "Hello ${NAME}!",
			want:    true,
		},
		{
			name:    "no variables",
			content: "Hello World!",
			want:    false,
		},
		{
			name:    "empty",
			content: "",
			want:    false,
		},
		{
			name:    "partial variable syntax",
			content: "Hello ${incomplete",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsAnyVariable(tt.content); got != tt.want {
				t.Errorf("ContainsAnyVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateVariables(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "all known variables",
			content: "${INPUT} ${DATE} ${DATETIME} ${CLIPBOARD} ${FILE}",
			want:    nil,
		},
		{
			name:    "unknown variable",
			content: "${UNKNOWN}",
			want:    []string{"UNKNOWN"},
		},
		{
			name:    "mixed known and unknown",
			content: "${INPUT} ${CUSTOM}",
			want:    []string{"CUSTOM"},
		},
		{
			name:    "no variables",
			content: "Plain text",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateVariables(tt.content)
			if len(got) != len(tt.want) {
				t.Errorf("ValidateVariables() = %v, want %v", got, tt.want)
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("ValidateVariables()[%d] = %q, want %q", i, v, tt.want[i])
				}
			}
		})
	}
}
