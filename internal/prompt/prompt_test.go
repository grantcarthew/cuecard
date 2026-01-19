package prompt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateFilename(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		existingFiles []string
		want          string
	}{
		{
			name:          "simple title",
			title:         "Hello World",
			existingFiles: nil,
			want:          "hello-world.md",
		},
		{
			name:          "title with special chars",
			title:         "What's This? A Test!",
			existingFiles: nil,
			want:          "whats-this-a-test.md",
		},
		{
			name:          "title with numbers",
			title:         "Prompt 123",
			existingFiles: nil,
			want:          "prompt-123.md",
		},
		{
			name:          "duplicate filename",
			title:         "Hello World",
			existingFiles: []string{"hello-world.md"},
			want:          "hello-world-2.md",
		},
		{
			name:          "multiple duplicates",
			title:         "Test",
			existingFiles: []string{"test.md", "test-2.md"},
			want:          "test-3.md",
		},
		{
			name:          "empty title",
			title:         "",
			existingFiles: nil,
			want:          "prompt.md",
		},
		{
			name:          "only special chars",
			title:         "!@#$%",
			existingFiles: nil,
			want:          "prompt.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateFilename(tt.title, tt.existingFiles); got != tt.want {
				t.Errorf("GenerateFilename() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	prompts := []*Prompt{
		{Title: "Code Review", Group: "Coding", Tags: []string{"review"}},
		{Title: "Bug Fix", Group: "Coding", Tags: []string{"fix", "debug"}},
		{Title: "Write Email", Group: "Writing", Description: "Draft professional emails"},
		{Title: "Refactor", Description: "Clean up code"},
	}

	tests := []struct {
		name  string
		query string
		want  int
	}{
		{"empty query returns all", "", 4},
		{"filter by title", "code", 2},              // "Code Review" (title) and "Refactor" (description has "code")
		{"filter by group", "coding", 2},            // Both in Coding group
		{"filter by tag", "review", 1},              // Only "Code Review" has review tag
		{"filter by description", "email", 1},       // Only "Write Email"
		{"case insensitive", "CODE", 2},             // Same as "code"
		{"no matches", "nonexistent", 0},            // No matches
		{"partial match", "fix", 1},                 // Only "Bug Fix" matches
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(prompts, tt.query)
			if len(got) != tt.want {
				t.Errorf("Filter() returned %d prompts, want %d", len(got), tt.want)
			}
		})
	}
}

func TestGroupPrompts(t *testing.T) {
	prompts := []*Prompt{
		{Title: "A", Group: "Group1", Favorite: true},
		{Title: "B", Group: "Group1"},
		{Title: "C", Group: "Group2"},
		{Title: "D", Favorite: true},
		{Title: "E"}, // ungrouped
	}

	favorites, groups, ungrouped := GroupPrompts(prompts)

	if len(favorites) != 2 {
		t.Errorf("favorites = %d, want 2", len(favorites))
	}

	if len(groups) != 2 {
		t.Errorf("groups = %d, want 2", len(groups))
	}

	if len(groups["Group1"]) != 2 {
		t.Errorf("Group1 = %d, want 2", len(groups["Group1"]))
	}

	if len(groups["Group2"]) != 1 {
		t.Errorf("Group2 = %d, want 1", len(groups["Group2"]))
	}

	if len(ungrouped) != 2 { // D and E are ungrouped
		t.Errorf("ungrouped = %d, want 2", len(ungrouped))
	}
}

func TestSortedGroupNames(t *testing.T) {
	groups := map[string][]*Prompt{
		"Zebra":  {},
		"Apple":  {},
		"Middle": {},
	}

	names := SortedGroupNames(groups)

	if len(names) != 3 {
		t.Fatalf("got %d names, want 3", len(names))
	}

	if names[0] != "Apple" || names[1] != "Middle" || names[2] != "Zebra" {
		t.Errorf("names = %v, want [Apple Middle Zebra]", names)
	}
}

func TestColorIndexForGroup(t *testing.T) {
	tests := []struct {
		name  string
		group string
		want  int
	}{
		{"empty group", "", -1},
		{"consistent hash", "Coding", ColorIndexForGroup("Coding")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ColorIndexForGroup(tt.group)
			if got != tt.want {
				t.Errorf("ColorIndexForGroup(%q) = %d, want %d", tt.group, got, tt.want)
			}
		})
	}

	// Verify hash is in valid range
	for _, group := range []string{"A", "Test", "Coding", "Writing", "Long Group Name"} {
		idx := ColorIndexForGroup(group)
		if idx < 0 || idx > 7 {
			t.Errorf("ColorIndexForGroup(%q) = %d, want 0-7", group, idx)
		}
	}
}

func TestPrompt_Validate(t *testing.T) {
	tests := []struct {
		name       string
		prompt     Prompt
		wantErrors int
		wantWarns  int
	}{
		{
			name: "valid prompt",
			prompt: Prompt{
				Title:       "Test",
				Description: "A test prompt",
			},
			wantErrors: 0,
			wantWarns:  0,
		},
		{
			name: "missing title",
			prompt: Prompt{
				Description: "No title",
			},
			wantErrors: 1,
			wantWarns:  0,
		},
		{
			name: "missing description",
			prompt: Prompt{
				Title: "No description",
			},
			wantErrors: 0,
			wantWarns:  1,
		},
		{
			name: "invalid input value",
			prompt: Prompt{
				Title:       "Test",
				Description: "Test",
				Input:       "invalid",
			},
			wantErrors: 1,
			wantWarns:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := tt.prompt.Validate()
			errors := 0
			warns := 0
			for _, r := range results {
				if r.Level == ValidationError {
					errors++
				} else {
					warns++
				}
			}
			if errors != tt.wantErrors {
				t.Errorf("Validate() errors = %d, want %d", errors, tt.wantErrors)
			}
			if warns != tt.wantWarns {
				t.Errorf("Validate() warnings = %d, want %d", warns, tt.wantWarns)
			}
		})
	}
}

func TestPrompt_ToMarkdown(t *testing.T) {
	p := &Prompt{
		Title:       "Test Prompt",
		Description: "A test",
		Group:       "Testing",
		Tags:        []string{"test", "example"},
		Input:       "required",
		InputHint:   "Enter text",
		Content:     "Hello ${INPUT}!",
	}

	md := p.ToMarkdown()

	// Verify it can be parsed back
	parsed, err := Parse(md)
	if err != nil {
		t.Fatalf("failed to parse generated markdown: %v", err)
	}

	if parsed.Title != p.Title {
		t.Errorf("round-trip Title = %q, want %q", parsed.Title, p.Title)
	}
	if parsed.Description != p.Description {
		t.Errorf("round-trip Description = %q, want %q", parsed.Description, p.Description)
	}
	if parsed.Group != p.Group {
		t.Errorf("round-trip Group = %q, want %q", parsed.Group, p.Group)
	}
	if parsed.Input != p.Input {
		t.Errorf("round-trip Input = %q, want %q", parsed.Input, p.Input)
	}
	if parsed.InputHint != p.InputHint {
		t.Errorf("round-trip InputHint = %q, want %q", parsed.InputHint, p.InputHint)
	}
	if parsed.Content != p.Content {
		t.Errorf("round-trip Content = %q, want %q", parsed.Content, p.Content)
	}
}

func TestLoadDirectory(t *testing.T) {
	// Create temp directory with test files
	tempDir := t.TempDir()

	// Create valid prompt
	validContent := `---
title: Valid Prompt
---

Content here.`
	if err := os.WriteFile(filepath.Join(tempDir, "valid.md"), []byte(validContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create README (should be skipped)
	if err := os.WriteFile(filepath.Join(tempDir, "README.md"), []byte("# Readme"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create non-markdown file (should be skipped)
	if err := os.WriteFile(filepath.Join(tempDir, "notes.txt"), []byte("notes"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create another valid prompt
	anotherContent := `---
title: Another Prompt
group: Testing
---

More content.`
	if err := os.WriteFile(filepath.Join(tempDir, "another.md"), []byte(anotherContent), 0644); err != nil {
		t.Fatal(err)
	}

	prompts, err := LoadDirectory(tempDir)
	if err != nil {
		t.Fatalf("LoadDirectory() error = %v", err)
	}

	if len(prompts) != 2 {
		t.Errorf("LoadDirectory() returned %d prompts, want 2", len(prompts))
	}
}

func TestCreatePromptFile(t *testing.T) {
	tempDir := t.TempDir()

	p := &Prompt{
		Title:   "New Prompt",
		Group:   "Test",
		Content: "Test content",
	}

	path, err := CreatePromptFile(tempDir, p)
	if err != nil {
		t.Fatalf("CreatePromptFile() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("file was not created")
	}

	// Verify content
	loaded, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}

	if loaded.Title != p.Title {
		t.Errorf("loaded Title = %q, want %q", loaded.Title, p.Title)
	}
}

func TestPrompt_GetVariables(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "single variable",
			content: "Hello ${INPUT}!",
			want:    []string{"INPUT"},
		},
		{
			name:    "multiple variables",
			content: "${INPUT} on ${DATE}",
			want:    []string{"INPUT", "DATE"},
		},
		{
			name:    "no variables",
			content: "Plain text",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prompt{Content: tt.content}
			got := p.GetVariables()
			if len(got) != len(tt.want) {
				t.Errorf("GetVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}
